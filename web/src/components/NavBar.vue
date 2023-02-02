<template>
    <nav class="navbar navbar-expand-lg navbar-dark bg-dark">
        <div class="container">
            <router-link class="navbar-brand" :to="{ name: 'home' }">GIN OJ</router-link>
            <div class="collapse navbar-collapse" id="navbarNavAltMarkup">
                <div class="navbar-nav me-auto mb-2 mb-lg-0">
                    <router-link :class="route_name == 'problem' ? 'nav-link active' : 'nav-link'"
                        :to="{ name: 'problem' }">问题列表</router-link>
                    <router-link :class="route_name == 'submit' ? 'nav-link active' : 'nav-link'"
                        :to="{ name: 'submit' }">我的提交</router-link>
                    <router-link :class="route_name == 'ranklist' ? 'nav-link active' : 'nav-link'"
                        :to="{ name: 'ranklist' }">排行榜</router-link>
                </div>
                <div class="navbar-nav" v-if="$store.state.user.is_login">
                    <li class="nav-item dropdown">
                        <a class="nav-link dropdown-toggle" href="#" role="button" data-bs-toggle="dropdown"
                            aria-expanded="false">
                            {{ $store.state.user.name }}

                        </a>
                        <ul class="dropdown-menu">
                            <li><router-link class="dropdown-item" :to="{ name: 'user-space' }">个人中心</router-link></li>
                            <li>
                                <hr class="dropdown-divider">
                            </li>
                            <li><a class="dropdown-item" href="#" @click="logout">退出</a></li>
                        </ul>
                    </li>

                </div>
                <div class="navbar-nav" v-else-if="!$store.state.user.pulling_info">
                    <li class="nav-item dropdown">
                        <router-link class="nav-link" :to="{ name: 'login' }" role="button">
                            登录
                        </router-link>
                    </li>
                    <li class="nav-item dropdown">
                        <router-link class="nav-link" :to="{ name: 'register' }" role="button">
                            注册
                        </router-link>
                    </li>

                </div>

            </div>
        </div>
    </nav>

</template>

<script>
import { useRoute } from 'vue-router';
import { computed } from 'vue'
import { useStore } from 'vuex';


export default {
    setup() {
        const route = useRoute();
        const store = useStore();
        let route_name = computed(() => route.name);

        const logout = () => {
            store.dispatch("logout")
        }
        return {
            route_name,
            logout,
        };

    },
}
</script>

<style scoped>

</style>