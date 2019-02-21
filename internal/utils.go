package internal

import (
    "bufio"
    "encoding/json"
    "fmt"
    "github.com/emacampolo/gomparator/internal/http"
    "github.com/urfave/cli"
    "io"
    "log"
    "os"
    "strconv"
    "strings"
)

func ReadFile(file *os.File) <-chan string {
    out := make(chan string)
    scanner := bufio.NewScanner(file)
    go func() {
        for scanner.Scan() {
            out <- scanner.Text()
        }
        close(out)
    }()

    return out
}

func DecodeLines(in <-chan string, headers map[string]string) <-chan string {
    out := make(chan string)
    go func() {
        for line := range in {
            line, err := decodeLine(line, headers)
            if err == nil {
                out <- line
            }
        }
        close(out)
    }()
    return out

}

func decodeLine(line string, headers map[string]string) (string, error) {
    if !strings.Contains(line, "access_token") {
        return line, nil
    }

    split := strings.Split(line, "-")
    if len(split) == 1 {
        return line, nil
    }

    clientID := split[len(split)-1]
    _, err := strconv.Atoi(split[1])
    if err != nil {
        return line, nil
    }

    response, err := http.Get(fmt.Sprintf("http://api.internal.ml.com/users/%s", clientID), headers)
    if err != nil {
        return "", nil
    }

    if response.StatusCode != 200 {
        return "", nil
    }

    var user map[string]interface{}
    if err := json.Unmarshal(response.Body, &user); err != nil {
        return "", nil
    }

    appID := split[1]
    response, err = http.Get(fmt.Sprintf("http://api.internal.ml.com/applications/%s/credentials?caller.id=%s&owner_id=%s&site_id=%s",
        appID, clientID, clientID, user["site_id"]), headers)
    if err != nil {
        return "", nil
    }

    if response.StatusCode != 200 {
        return "", nil
    }

    var token map[string]interface{}
    if err := json.Unmarshal(response.Body, &token); err != nil {
        return "", nil
    }

    toReplace := "APP_USR-" + appID + "-" + split[2] + "-X-" + clientID
    return strings.Replace(line, toReplace, token["access_token"].(string), 1), nil
}

func ParseHeaders(c *cli.Context) map[string]string {
    var result map[string]string

    headers := strings.Split(c.String("headers"), ",")
    result = make(map[string]string, len(headers))

    for _, header := range headers {
        if header == "" {
            continue
        }

        h := strings.Split(header, ":")
        if len(h) != 2 {
            log.Fatal("invalid header")
        }

        result[h[0]] = h[1]
    }

    return result
}

func Close(c io.Closer) {
    err := c.Close()
    if err != nil {
        log.Fatal(err)
    }
}
