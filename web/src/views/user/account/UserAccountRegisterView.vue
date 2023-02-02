<template>

    <ContainerCard>
        <div class="row justify-content-md-center">
            <div class="col-3">
                <form @submit.prevent="login">
                    <div class="mb-3">
                        <label for="username" class="form-label">用户名</label>
                        <input v-model="username" type="text" class="form-control" id="username" placeholder="请输入用户名">
                    </div>
                    <div class="mb-3">
                        <label for="password" class="form-label">密码</label>
                        <input v-model="password" type="password" class="form-control" id="password"
                            placeholder="请输入密码">
                    </div>
                    <div class="mb-3">
                        <label for="email" class="form-label">邮箱</label>
                        <input v-model="mail" type="text" class="form-control" id="mail" placeholder="请输入邮箱">
                    </div>

                    <div class="mb-3">
                        <label for="code" class="form-label">验证码</label>
                        <input v-model="code" type="text" class="form-control" id="code" placeholder="请输入验证码">
                    </div>


                    <div class="error-message">{{ error_message }}</div>
                    <button type="submit" class="btn btn-primary send" @click="send_code">发送验证码</button>
                    <button type="submit" class="btn btn-primary register" @click="register">注册</button>

                </form>
            </div>
        </div>
    </ContainerCard>


</template>

<script>
import ContainerCard from "@/components/ContainerCard.vue"
import { ref } from "vue";
// import { useStore } from 'vuex'
import router from "@/router";//export default 不叫{}
import $ from "jquery";
export default {
    components: {
        ContainerCard,
    },

    setup() {
        // const store = useStore();
        let username = ref('');
        let password = ref('');
        let code = ref('');
        let mail = ref('');
        let error_message = ref('');
        const send_code = () => {
            $.ajax({
                type: "post",
                url: "http://127.0.0.1:3000/user-send_code",
                data: {
                    email: mail.value
                },
                success: function (response) {
                    if (response.code != "200") {
                        error_message.value = response.msg;
                    }

                },
                error(resp) {
                    error_message.value = resp.msg;
                }
            });
        }
        const register = () => {
            $.ajax({
                type: "post",
                url: "http://127.0.0.1:3000/user-register",
                data: {
                    name: username.value,
                    password: password.value,
                    code: code.value,
                    mail: mail.value
                },
                success: function (response) {
                    if (response.code == "200") {
                        router.push({ name: "login" });
                    } else {
                        error_message.value = response.msg;
                    }
                },
                error(resp) {
                    error_message.value = resp.msg;

                }
            });
        }
        return {
            username,
            password,
            code,
            mail,
            error_message,
            register,
            send_code,
        }
    }
}


</script>

<style scoped>
button.register {
    float: right;
    width: 40%;
}

button.send {
    width: 40%;
}

div.error-message {
    color: red;
}
</style>