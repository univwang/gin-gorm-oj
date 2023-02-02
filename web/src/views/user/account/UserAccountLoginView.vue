<template>
    <ContainerCard v-if="!$store.state.user.pulling_info">
        <div class="row justify-content-md-center">
            <div class="col-3">
                <form @submit.prevent="login">
                    <div class="mb-3">
                        <label for="username" class="form-label">用户名</label>
                        <input v-model="username" type="text" class="form-control" id="username" placeholder="请输入用户名">
                    </div>
                    <div class="mb-3">
                        <label for="password" class="form-label">请输入密码</label>
                        <input v-model="password" type="password" class="form-control" id="password"
                            placeholder="请输入密码">
                    </div>
                    <div class="error-message">{{ error_message }}</div>
                    <button type="submit" class="btn btn-primary">登录</button>

                </form>
            </div>
        </div>
    </ContainerCard>
</template>

<script>
import ContainerCard from "@/components/ContainerCard.vue"
import { ref } from "vue";
import { useStore } from "vuex";
import router from '@/router/index'

export default {
    components: {
        ContainerCard,
    },
    setup() {
        const store = useStore();
        let username = ref('');
        let password = ref('');
        let error_message = ref('');
        const token = localStorage.getItem("token");
        if (token) {
            store.commit("updateToken", token);
            store.dispatch("getinfo", {
                success(resp) {
                    if (resp.code == "200") {
                        router.push({ name: "home" })
                        store.commit("updatePullingInfo", false);
                    } else {
                        store.commit("updatePullingInfo", false);
                    }
                },
                error() {
                    store.commit("updatePullingInfo", false);
                }
            })
        } else {
            store.commit("updatePullingInfo", false);
        }
        const login = () => {
            store.dispatch("login", {
                username: username.value,
                password: password.value,
                success() {
                    store.dispatch("getinfo", {
                        success() {
                            router.push({ name: 'home' });
                            // console.log(store.state);
                        },
                        error(resp) {
                            console.log(resp);
                        }
                    })
                },
                error(resp) {
                    error_message.value = resp.msg
                }
            })
        }

        return {
            username,
            password,
            error_message,
            login,
        }
    }
}
</script>

<style scoped>
button {
    width: 100%;
}

div.error-message {
    color: red;
}
</style>