# APIs
과학미科學美의 API 명세

## Account API
### Login API
- `POST api/account/login`
  - request header: X
  - params: X
  - request body:
    - uid(String): 아이디 값인 이메일
    - pw(String): 사용자 비밀번호
  - response header: X
  - response body:
    - uname(String): 사용자 이름
    - isSuccess(Boolean): 로그인 성공 여부
    - message(String): 로그인 결과에 대한 상세한 메시지
      - ex) `잘못된 PW`, `로그인 성공` 등
### Signup API
- `POST api/account/signup`
  - request header: X
  - params: X
  - request body:
    - uname(String): 사용자 이름
    - uid(String): 아이디 값인 이메일
    - pw(String): 사용자 비밀번호
  - response header: X
  - response body:
    - uname(String): 사용자 이름
    - isSuccess(Boolean): signup 성공 여부
    - message(String): signup 결과에 대한 상세한 메시지
      - ex) `이미 존재하는 계정`, `Signup 성공` 등
