# gRPC Deep Dive
* Use [Protocol Buffer](https://developers.google.com/protocol-buffers) for communication
* Why Protocol Buffer?
    * message definition 을 작성하기 쉽다.
    * 간단한 .proto 파일 작성으로 여러 언어의 많은 양의 코드들을 만들기 쉽다.
    * 페이로드는 binary이고 그러므로 네트워크로 받고 serialize, de-serialize가 효율적이다.
    * 마이크로서비스를 만들기 좋다.
* HTTP2
    * http1.1의 문제?
        * 요청할때마다 새로운 TCP 연결을 맺는다.
        * 압축되지 않은 헤더를 포함하기 때문에 무겁다.
        * request/response 만 가능하다 (서버 푸쉬가 되지 않음)
        * GET/POST 뿐이다.
    * Multiplexing 지원
    * 하나의 TCP 연결에서 클라이언트와 서버가 병렬적으로 푸쉬 할 수 있다.
        * 위 결과 latency를 크게 줄인다.
    * 서버 푸쉬를 지원한다.
    * 클라이언트가 한번 요청시에도 서버가 여러번 푸쉬 할 수 있다. 
    * 헤더 압축을 지원한다.
        * 패킷 사이즈를 줄인다. 
    * Binary이다
        * 네트워크 전송에 효율적이고 Protocol Buffer (Binary Protocol)과 잘 맞는다.
    * 보안성이 높다.
* Type of gRPC
    * Unary
        * 기존 방식과 같이 한번 요청하면 한번 응답
    * Server Streaming
        * 클라이언트가 한번 요청하면 서버가 스트리밍 처럼 여러번 응답한다.
    * Client Streaming
        * 클라이언트가 스트리밍하게 보내고 마지막에 서버가 응답
    * Bi Directional Streaming
        * 양방향 스트리밍
* Scalability in gRPC
    * gRPC 서버는 기본값이 비동기식이다.
        * 요청시 스레드를 block하지 않는다.
        * 서버가 병렬적으로 요청을 serve할 수 있다. 
    * gRPC 클라이언트는 동기적이거나 비동기적일 수 있다. 
        * 클라이언트가 성능의 요구에 따라 결정할 수 있다. 
    * gRPC Client 는 client side load balancing을 할 수 있다.
* gRPC VS REST
|gRPC|REST|
|---|---|
|Protocol Buffer (small, fast)|JSON (text based, slow, big)|
|HTTP/2 (lower latency)|HTTP/1.1 (higher latency)|
|Bidirectional & Async|Client => Server requests only|
|Stream support|Request / Response support only|
|API Oriented - “What”|CRUD Oriented|
|Code Generation through Protocol Buffers in any language (1st class citizen)|Code generation through OpenAPI / Swagger (add-on) (2nd class citizen)|
|RPC Based|HTTP verbs based|

