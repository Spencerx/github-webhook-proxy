package main

import (
  "bytes"
  "encoding/json"
  "io/ioutil"
  "log"
  "net/http"
  "time"
)

type WebhookProxyClient struct {
}

func (c *WebhookProxyClient) request() {
  client := &http.Client{}
  req, err := http.NewRequest(method, endpoint, reader)
  if err != nil {
    return err
  }
  req.Header = headers

  resp, err := client.Do(req)
  if err != nil {
    return nil, err
  }

  defer resp.Body.Close()
  if resp.StatusCode >= 400:
    return nil, errors.New("Not found")
  }
}

func (c *WebhookProxyClient) Run() {
  log.Printf("Forwarding webhooks from %s to %s", proxyEndpoint, jenkinsEndpoint)

  ticker := time.NewTicker(10 * time.Second)
  for _ = range ticker.C {
    resp, err := http.Get(proxyEndpoint + "/webhook")
    if err != nil {
      log.Printf("Failed to retrieve webhooks: %s", err.Error())
      continue
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
      log.Printf("Failed to read query: %s", err.Error())
      resp.Body.Close()
      continue
    }

    var requests RequestCache
    if err := json.Unmarshal(body, &requests.Requests); err != nil {
      log.Printf("Failed to unmarshal requests: %s: %s", err.Error(), string(body))
      resp.Body.Close()
      continue
    }

    resp.Body.Close()

    if len(requests.Requests) == 0 {
      continue
    }

    for _, request := range requests.Requests {
      log.Printf("Resending request %+v\n", request)

      req, err := c.request(request.Method, jenkinsEndpoint, bytes.NewReader([]byte(request.Body)))
      if err != nil {
        log.Printf("Failed to create request: %+s", err.Error())
        continue
      }
      req.Header = request.Headers

      resp, err := client.Do(req)
      if err != nil || resp.StatusCode >= 400 {
        log.Printf("Failed to push webhook: %+v %+v", err, resp)
        continue
      }
      resp.Body.Close()

      req, err = http.NewRequest("DELETE", proxyEndpoint+"/webhook/"+request.UUID, nil)
      if err != nil {
        log.Printf("Failed to create request: %+s", err.Error())
        continue
      }

      resp, err = client.Do(req)
      if err != nil || resp.StatusCode >= 400 {
        log.Printf("Failed to delete webhook %s: %+v, %+v", request.UUID, err, resp)
      }

      if resp != nil {
        resp.Body.Close()
      }
    }
  }
}

func NewWebhookProxyClient() *WebhookProxyClient {
  return &WebhookProxyClient{}
}
