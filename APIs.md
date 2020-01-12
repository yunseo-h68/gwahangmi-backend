# APIs
과학미科學美의 API 명세 및 Description

## APIs 목차
- [API Response 구조](#api-response-구조)
- [Account API](#account-api)
  - [Login API](#login-api)
  - [Signup API](#signup-api)
  - [Profile API](#profile-api)
  - [Users API](#users-api)
  - [User API](#user-api)
- [Category API](#category-api)
  - [Posts API](#posts-api)
  - [Post API](#post-api)
  - Point API
  - Comment API
- SciQuiz API
  - Quizzes API
  - Quiz API
- File API

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
### Users API
- `GET /api/account/users`
  - request header: X
  - params: 
    - limit(int): 조회할 글(Post)의 개수
    - point(Boolean): true일 경우, point순
    - post(Boolean): true일 경우, post의 개수순
    - sort(Boolean): true면 오름차순, false이면 내림차순 정렬
  - request body: X
  - response header: 
    - `Content-Type`: `application/json`
  - response body:
    - users([]String): 사용자의 ID 배열
### User API
- `GET /api/account/users/:uid`
  - request header: X
  - params: X
  - request body: X
  - response header: 
    - `Content-Type`: `application/json`
  - response body:
    - uid(String): 사용자 ID
    - uname(String): 사용자 이름
    - profileImg(String): 사용자 프로필 이미지 파일 이름
    - point(int): 사용자가 보유한 포인트
    - postCnt(int): 사용자가 게시한 글의 개수

## Category API
### Posts API
- `GET /api/category/posts`
  - request header: X
  - params: 
    - limit(int): 조회할 글(Post)의 개수
    - popularity(Boolean): true일 경우, 기준을 인기순으로 잡음. false일 경우 최신순
    - total(Boolean): popularity가 true이고 total도 true일 경우, point의 총점순
    - average(Boolean): popularity가 true이고 average도 true일 경우, point의 평균순
    - sort(Boolean): true면 오름차순, false이면 내림차순 정렬
  - request body: X
  - response header: 
    - `Content-Type`: `application/json`
  - response body:
    - posts([]String): 글(Post)의 ID 배열
- `POST /api/category/posts`
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
- `PUT /api/category/posts`
  - request header: 
    - `Content-Type`: `application/json`
  - params: X
  - request body:
    - author(String): 수정한 사용자의 ID
    - category(String): 수정한 글의 카테고리(분류)
    - postID(string): 수정할 글의 ID
    - title(String): 수정한 글의 제목
    - content(String): 수정한 글의 본문
  - response header: 
    - `Content-Type`: `application/json`
  - response body:
    - postID(String): 수정한 글의 ID
    - isSuccess(Boolean): 글 수정 성공 여부
    - message(String): 글 수정 결과에 대한 상세한 메시지
- `DELETE /api/category/posts`
  - request header: X
  - params:
    - post_id(String): 삭제할 글의 ID
    - category(String): 수정할 글의 카테고리(분류)
  - request body: X
  - response header: 
    - `Content-Type`: `application/json`
  - response body:
    - postID(String): 삭제 시도한 글의 ID(삭제 성공시 `""`)
    - isSuccess(Boolean): 글 삭제 성공 여부
    - message(String): 글 삭제 결과에 대한 상세한 메시지 
### Post API
- `GET /api/category/posts/:post-id`
  - request header: X
  - params: X
  - request body: X
  - response header: 
    - `Content-Type`: `application/json`
  - response body:
    - postID(String): 해당 글(Post)의 ID
    - author(String): 해당 글의 작성자 ID
    - category(String): 해당 글의 카테고리(분류)
    - title(String): 해당 글의 제목
    - content(String): 해당 글의 본문(내용) 파일 이름
    - participantCnt(int): 해당 글 평가에 참여한 사용자 수
    - totalPoint(int): 해당 글의 총 평가 점수(point)
    - averagePoint(float): 해당 글의 평균 평가 점수(point)
    - uploadDate:
      - year: 해당 글의 작성시간 중 연도
      - month: 해당 글의 작성시간 중 월
      - day: 해당 글의 작성시간 중 일
      - hour: 해당 글의 작성시간 중 시간
      - minute: 해당 글의 작성시간 중 분
      - second: 해당 글의 작성시간 중 초
      - fullDate(String): 해당 글의 작성시간