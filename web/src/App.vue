<template>
  <div>
    <div>用户名： {{ user_name }}</div>
    <div>邮箱: {{ mail }}</div>
  </div>
  <router-view></router-view>
</template>

<script>
import $ from 'jquery';
import { ref } from 'vue';
export default {
  name: "App",
  setup: () => {
    let user_name = ref("");
    let mail = ref("");

    $.ajax({
      type: "get",
      url: "http://127.0.0.1:3000/user-detail?identity=" + "b7faff77-b741-49a3-aa72-6cf0b953c72c",
      success: resp => {
        let data = JSON.parse(JSON.stringify(resp))
        if (data.code != "200") {
          console.log("获取用户信息错误")
        } else {
          user_name.value = data.msg.name;
          mail.value = data.msg.mail
        }
      }
    });
    return {
      user_name,
      mail
    };
  }
}
</script>

<style>
body {
  background-image: url(@/assets/background.jpg);
  background-size: cover;
}
</style>