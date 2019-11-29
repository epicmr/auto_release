<template>
    <div id="login" class="login">
        <el-row class="row">
            <img src="../assets/logo.png">
        </el-row>
        <el-row class="row">
            <el-input class="input" v-model="phone" placeholder="请输入手机号" clearable></el-input>
        </el-row>
        <el-row class="row">
            <el-input class="input" v-model="passwd" placeholder="请输入密码" clearable show-password></el-input>
        </el-row>
        <el-row class="row">
            <el-button class="button" type="primary" @click="login()">登陆</el-button>
        </el-row>
        <el-row class="row">
            <router-link to="/session/regist" id="regist">立即注册</router-link>
        </el-row>
    </div>
</template>

<script>
export default {
  name: 'v-login',
  data () {
    return {
      phone: '',
      passwd: ''
    }
  },
  methods: {
    login () {
      console.log('login')
        this.$http.post('/session/login', {'phone':this.phone,'passwd':this.passwd})
    }
  }
}
</script>

<style lang="scss">
.login {
    position: relative;
    text-align:center;
    margin: 0 auto;
    border-radius: 4px;
    font-size: 14px;
    align:center;
    .input {
        width: 300px;
    }
    .button {
        width: 120px;
    }
    .row {
        margin: 10px;
    }
}

.table {
    .demo-table-expand {
        font-size: 0;
    }

    .demo-table-expand label {
        width: 90px;
        color: #99a9bf;
    }

    .demo-table-expand .el-form-item {
        margin-right: 0;
        margin-bottom: 0;
        width: 30%;
    }
}

</style>
