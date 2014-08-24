wrk.method = "POST"
wrk.body = "grant_type=password&username=email@example.com&password=password"
wrk.headers["Content-Type"] = "application/x-www-form-urlencoded"
wrk.headers["Authorization"] = "Basic Y2xpZW50X2lkOmNsaWVudF9zZWNyZXQ="
