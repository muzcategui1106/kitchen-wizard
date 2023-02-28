package swagger

import "net/http"

const UIPrefix = "/swagger-ui/"

var Handler = http.FileServer(http.Dir("./swagger-ui"))
