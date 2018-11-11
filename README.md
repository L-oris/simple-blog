### Currently Implemented

-   Dependency Injection with `https://github.com/sarulabs/di`
    -   then replaced by Go Cloud's compile-time DI
        -   https://blog.golang.org/wire
-   DB connection to MySQL (still locally)
-   Env variables > `.env` file
-   Independent Controllers (routing) injected in main Router
    -   only partially tested
-   Google Cloud Storage for image uploading/downloading
    -   `https://github.com/google/go-cloud`
-   Server side rendering with Go Templates

### Todo

-   separate collections for users & posts
-   auth with cookies and CSRF
-   deploy to Google Cloud, also DB
-   docker
-   circleCI

