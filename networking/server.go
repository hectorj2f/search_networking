package networking

import (
  "crypto/tls"
  "io"
  "net"
  "fmt"
  "strconv"

  "github.com/hectorj2f/search_networking/api"
  "github.com/hectorj2f/search_networking/resources"
  "github.com/hectorj2f/search_networking/networking/serializer"
  "github.com/hectorj2f/search_networking/types"

  "github.com/docker/libchan"
  "github.com/docker/libchan/spdy"
  logger "github.com/Sirupsen/logrus"
)

func startListener(cert string, key string, port int) (net.Listener, error){
  var listener net.Listener

  if cert != "" && key != "" {
    tlsCert, err := tls.LoadX509KeyPair(cert, key)
    if err != nil {
      return nil, err
    }

    tlsConfig := &tls.Config{
      InsecureSkipVerify: true,
      Certificates:       []tls.Certificate{tlsCert},
    }

    listener, err = tls.Listen("tcp", fmt.Sprintf("%s:%s", resources.SERVER_ADDR, strconv.Itoa(port)), tlsConfig)
    if err != nil {
      return nil, err
    }
  } else {
    var err error
    logger.Info("TLS configuration is disabled")
    listener, err = net.Listen("tcp", fmt.Sprintf("%s:%s", resources.SERVER_ADDR, strconv.Itoa(port)))
    if err != nil {
      return nil, err
    }
  }
  logger.Infof("Listening at %s...\n", strconv.Itoa(port))
  return listener, nil
}


func Server(cert string, key string, port int) error {
  listener, err := startListener(cert, key, port)
  if err != nil {
    logger.Error(err)
    return err
  }
  defer listener.Close()

  tl, err := spdy.NewTransportListener(listener, spdy.NoAuthenticator)
  if err != nil {
    logger.Error(err)
    return err
  }

  for {
    transport, err := tl.AcceptTransport()
    if err != nil {
      return err
    }

    logger.Info("Waiting for receive channel...")
    receiver, err := transport.WaitReceiveChannel()
    if err != nil {
      return err
    }

    // Go routine to wait for incoming messages
    go waitForMessages(receiver, transport)

  }
  return nil
}

func waitForMessages(receiver libchan.Receiver, transport *spdy.Transport) error {
  for {
      logger.Info("Waiting for connection...")
      message := &types.Message{}
      err := receiver.Receive(message)
      if err != nil && err != io.EOF {
        logger.Errorf("Error receiving message: %s", err)
      }

      if err == io.EOF {
        receiver, err = transport.WaitReceiveChannel()
        if err != nil {
          logger.Error(err)
        }
      } else {
        logger.Debug("Message received!")
        logger.Debug(message)

        if len(message.Data) < 1 {
          logger.Errorf("Unexpected message length: %s", strconv.Itoa(len(message.Data)))
        }

        query := serializer.DecodeMsgPack(message.Data)
        users, err := api.SearchUsers(query)
        if err != nil {
          logger.Errorf("Error searching the users!", err)
        }
        data, err := serializer.EncodeArrayMsgPack(users)

        ack := &types.AckMessage{Id: message.Id, Data: data}
        logger.Debug("Sending ack...")
        err = message.Ret.Send(ack)
        logger.Debug("Sent ack!")
        if err != nil {
          logger.Errorf("Error sending ack: %s", err)
        }
      }
  }

}
