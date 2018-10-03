### Currently Implemented

-   Dependency Injection with `https://github.com/sarulabs/di`
-   DB connection to MySQL (still locally)
-   Env variables > `.env` file
-   Independent Controllers (routing) injected in main Router
    -   only partially tested
-   Google Cloud Storage for image uploading/downloading
    -   `https://github.com/google/go-cloud` library, works great but still not safe for production
-   Server side rendering with Go Templates

### Todo

-   separate collections for users & posts
-   auth with cookies and CSRF
-   deploy to Google Cloud, also DB
-   docker
-   circleCI
