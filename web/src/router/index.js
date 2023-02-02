import { createRouter, createWebHistory } from 'vue-router'

import ProblemIndexView from "@/views/user/problem/ProblemIndexView"
import RankListIndexView from "@/views/user/ranklist/RanklistIndexView"
import SubmitIndexView from "@/views/user/problem/SubmitIndexView"
import NotFound from "@/components/NotFound"
import UserAccountLoginView from "@/views/user/account/UserAccountLoginView"
import UserAccountRegisterView from "@/views/user/account/UserAccountRegisterView"
import UserSpaceIndexView from "@/views/user/space/UserSpaceIndexView"
import store from '@/store'

const routes = [
  {
    path: "/",
    component: ProblemIndexView,
    redirect: '/problem/',
    name: 'home',
    meta: {
      requestAuth: true,
    }
  },
  {
    path: "/problem/",
    component: ProblemIndexView,
    name: 'problem',
    meta: {
      requestAuth: true,
    }
  },
  {
    path: "/user/space/",
    component: UserSpaceIndexView,
    name: 'user-space',
    meta: {
      requestAuth: true,
    }
  },
  {
    path: "/submit/",
    component: SubmitIndexView,
    name: 'submit',
    meta: {
      requestAuth: true,
    }
  },
  {
    path: "/user/account/login/",
    component: UserAccountLoginView,
    name: 'login'
  },
  {
    path: "/user/account/register/",
    component: UserAccountRegisterView,
    name: 'register'
  },
  {
    path: "/ranklist/",
    component: RankListIndexView,
    name: 'ranklist',
    meta: {
      requestAuth: true,
    }
  },
  {
    path: "/:catchAll(.*)",
    component: NotFound,
    name: '404'
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  if (to.meta.requestAuth && !store.state.user.is_login) {
    next({ name: "login" })
  } else {
    next()
  }
})

export default router
