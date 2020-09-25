/**
 * Created by sean on 2020/6/19.
 */
function modulesGet(username, success, failed) {
    let url = 'http://localhost:9090/api/v1/module/getall';
    let data = {userName: username};
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
