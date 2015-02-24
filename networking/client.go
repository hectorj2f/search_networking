package networking

import (
  "crypto/tls"
  "errors"
  "fmt"
  "net"
  "strconv"

  "github.com/hectorj2f/search_networking/networking/serializer"
  "github.com/hectorj2f/search_networking/types"

  "code.google.com/p/go-uuid/uuid"

  "github.com/docker/libchan"
  "github.com/docker/libchan/spdy"
  logger "github.com/Sirupsen/logrus"
)


func Query(query map[string]interface{}, server_addr string, use_tls bool, port int) ([]map[string]interface{}, error) {
  client, err := setupTLSConnection(server_addr, use_tls, port)
  if err != nil {
    return nil, err
  }
  transport, err := spdy.NewClientTransport(client)
  if err != nil {
    return nil, err
  }

  logger.Debug(query)
  encoded_data, err := serializer.EncodeMsgPack(query)

  sender, err := transport.NewSendChannel()
  if err != nil {
    logger.Error(err)
    return nil, err
  }
  defer sender.Close()

  receiver, remoteSender := libchan.Pipe()
  message, err := sendMessage(encoded_data, sender, remoteSender)

  ack, err := waitAckMessage(receiver, message)
  if err != nil {

  }
  result := serializer.DecodeArrayMsgPack(ack.Data)
  logger.Debug("Closing sender...")

  return result, nil
}

func setupTLSConnection(server_addr string, use_tls bool, port int) (net.Conn, error){
  var (
    client net.Conn
    err error
  )
  if use_tls {
    client, err = tls.Dial("tcp",
                            fmt.Sprintf("%s:%s", server_addr, strconv.Itoa(port)),
                            &tls.Config{InsecureSkipVerify: true})
  } else {
    client, err = net.Dial("tcp", fmt.Sprintf("%s:%s", server_addr, strconv.Itoa(port)))
  }

  return client, err
}

func sendMessage(encoded_data []byte, sender libchan.Sender, remoteSender libchan.Sender) (*types.Message, error){
  message := &types.Message{
      Data:    encoded_data,
      Id:      uuid.NewRandom(),
      Ret:     remoteSender,
  }
  logger.Debug("Sending a message ....")
  err := sender.Send(message)
  if err != nil {
    logger.Errorf("Error sending message: %s", err)
    return nil, err
  }

  return message, nil
}

func waitAckMessage(receiver libchan.Receiver, message *types.Message) (*types.AckMessage, error){
  logger.Debug("Waiting to receive ack...")
  ack := &types.AckMessage{}
  if err := receiver.Receive(ack); err != nil {
    logger.Errorf("Error receiving ack!")
    return nil, err
  }
  logger.Debug("Received ack!")

  if ack.Id.String() != message.Id.String() {
    logger.Errorf("Unexpected ACK identifier: %d", ack.Id)
    return nil, errors.New(fmt.Sprintf("Unexpected ACK identifier: %d", ack.Id))
  }

  return ack, nil
}
