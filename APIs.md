# APIs
과학미科學美의 API 명세 및 Description

## API Response 구조
- code(int) : 응답 코드
- message(String) : 에러메시지
- data : 해당 API의 response body

&lt;Example&gt;
- Login API Response Body 예시
```json
{
    "code": 200,
    "message": "",
    "data": {
        "uname": "Yunseo Hwang",
        "isSuccess": true,
        "message": "로그인에 성공하셨습니다"
    }
}
```

- Signup API Response Body 예시
```json
{
    "code": 200,
    "message": "",
    "data": {
        "uname": "",
        "isSuccess": false,
        "message": "계정이 이미 존재합니다"
    }
}
```

## Account API
### Login API
- `POST /api/account/login`
  - request header: 
    - `Content-Type`: `application/json`
  - params: X
  - request body:
    - uid(String): 아이디 값인 이메일
    - pw(String): 사용자 비밀번호
  - response header: 
    - `Content-Type`: `application/json`
  - response body:
    - uname(String): 사용자 이름
    - isSuccess(Boolean): 로그인 성공 여부
    - message(String): 로그인 결과에 대한 상세한 메시지
      - ex) `잘못된 PW`, `로그인 성공` 등
### Signup API
- `POST /api/account/signup`
  - request header: 
    - `Content-Type`: `application/json`
  - params: X
  - request body:
    - uname(String): 사용자 이름
    - uid(String): 아이디 값인 이메일
    - pw(String): 사용자 비밀번호
  - response header: 
    - `Content-Type`: `application/json`
  - response body:
    - uname(String): 사용자 이름
    - isSuccess(Boolean): signup 성공 여부
    - message(String): signup 결과에 대한 상세한 메시지
      - ex) `이미 존재하는 계정`, `Signup 성공` 등
### Profile API
- `GET /api/account/profile`
  - request header:
  - params: 
    - uid(String): 아이디 값인 이메일 
  - request body: X
  - response header: 
    - `Content-Type`: `application/json`
  - response body:
    - profile_img(String): 사용자 프로필의 프로필 이미지 경로
    - isSuccess(Boolean): 프로필 조회 성공 여부
    - message(String): 프로필 조회 결과에 대한 상세한 메시지
- `POST /api/account/profile`
  - request header:
    - `Content-Type`: `multipart/form-data`
  - params: X
  - request body:
    - uid(String) : 아이디 값인 이메일
    - profile_img(Binary Data): 프로필 이미지 파일
  - response header:
    - `Content-Type`: `application/json`
  - response body:
    - profile_img(String): 프로필 이미지 파일 경로
    - isSuccess(Boolean): 프로필 등록 성공 여부
    - message(String): 프로필 조회 결과에 대한 상세한 메시지
- `PUT /api/account/profile`
   - request header:
    - `Content-Type`: `multipart/form-data`
  - params: X
  - request body:
    - uid(String) : 아이디 값인 이메일
    - profile_img(Binary Data): 프로필 이미지 파일
  - response header:
    - `Content-Type`: `application/json`
  - response body:
    - profile_img(String): 프로필 이미지 파일 경로
    - isSuccess(Boolean): 프로필 수정 성공 여부
    - message(String): 프로필 조회 결과에 대한 상세한 메시지
- `DELETE /api/account/profile`
  - request header: X
  - params: 
    - uid(String): 아이디 값인 이메일 
  - request body: X
  - response header: 
    - `Content-Type`: `application/json`
  - response body:
    - profile_img(String): 프로필 삭제 후, Default 프로필 이미지 파일 경로
    - isSuccess(Boolean): 프로필 조회 성공 여부
    - message(String): 프로필 조회 결과에 대한 상세한 메시지

## Category API
### Posts API
- `GET /api/category/posts`
  - request header: X
  - params: 
    - limit(int): 조회할 글(Post)의 개수
  - request body: X
  - response header: 
    - `Content-Type`: `application/json`
  - response body:
    - postID(String): 글(Post)의 ID
    - isSuccess(Boolean): 글 조회 성공 여부
    - message(String): 글 조회 결과에 대한 상세한 메시지
- `POST api/category/posts`
  - request header: 
    - `Content-Type`: `application/json`
  - params: X
  - request body:
    - author(String): 작성한 사용자의 ID
    - category(String): 작성한 글의 카테고리(분류)
    - title(String): 작성한 글의 제목
    - content(String): 작성한 글의 본문
  - response header: 
    - `Content-Type`: `application/json`
  - response body:
    - postID(String): 작성한 글의 ID
    - isSuccess(Boolean): 글 업로드 성공 여부
    - message(String): 글 업로드 결과에 대한 상세한 메시지