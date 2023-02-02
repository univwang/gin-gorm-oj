import $ from "jquery";
export default {
    state: {
        identity: "",
        name: "",
        mail: "",
        photo: "",
        token: "",
        is_login: false,
        is_admin: 0,
        pulling_info: true, //是否在拉取信息
    },
    getters: {
    },
    mutations: {
        updatePullingInfo(state, pulling_info) {
            state.pulling_info = pulling_info;
        },
        updateUser(state, user) {
            state.identity = user.identity;
            state.name = user.name;
            state.mail = user.mail;
            state.photo = user.photo
            state.is_login = user.is_login
        },
        updateToken: (state, token) => {
            state.token = token
        },
        logout(state) {
            state.identity = "";
            state.name = "";
            state.photo = "";
            state.token = "";
            state.is_login = false;
        }
    },
    actions: {
        login(context, data) {
            $.ajax({
                url: "http://127.0.0.1:3000/user-login",
                type: "post",
                data: {
                    username: data.username,
                    password: data.password,
                },
                success(resp) {
                    // let jsonData = JSON.parse(JSON.stringify(resp));
                    let jsonData = resp
                    if (jsonData.code == "200") {
                        localStorage.setItem("token", jsonData.data.token);
                        context.commit("updateToken", jsonData.data.token) //更新数据
                        data.success(resp) //返回成功
                    } else {
                        data.error(resp) // 返回失败
                    }
                },
                error(resp) {
                    data.error(resp); //返回失败
                }
            })
        },
        getinfo(context, data) {
            $.ajax({
                url: "http://127.0.0.1:3000/me/user-info",
                type: "get",
                headers: {
                    Authorization: context.state.token,
                },
                success(resp) {
                    if (resp.code == "200") {
                        context.commit("updateUser", {
                            ...resp.msg,
                            is_login: true,
                        });
                        data.success(resp); //返回成功                        
                    } else {
                        data.error(resp); // 返回失败
                    }
                },
                error(resp) {
                    data.error(resp); //返回失败
                }
            })
        },
        logout(context) {
            localStorage.removeItem("token")
            context.commit("logout");
        }
    },
    modules: {
    }
}
