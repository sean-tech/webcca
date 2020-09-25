

function userLogin(username, password, success, failed) {
    let url = 'http://localhost:9090/api/v1/admin/login';
    let data = {userName: username,password:password,uuid:"azczxcz",client:"chrome"};
    let content = {
        method: 'POST', // or 'PUT'
        body: JSON.stringify(data),
        headers: new Headers({
            'Access-Control-Allow-Origin': '*',
            'Content-Type': 'application/json'
        })
    };

    fetch(url, content).then(res => res.json())
        .catch(error => failed(error))
        .then(response => {
            if (response.code != 200) {
                failed(response.msg)
            } else {
                success(response)
            }
        });
}



// {
//   "code": 0
//   ,"msg": "登入成功"
//   ,"data": {
//     "access_token": "c262e61cd13ad99fc650e6908c7e5e65b63d2f32185ecfed6b801ee3fbdd5c0a"
//   }
// }