<template>
    <div>
        <el-col id="header" class="site-header">
            <el-button class="logout-button" size="mini" @click="logout">注销</el-button>
        </el-col div>
    </div>
</template>

<script>
export default {
    name: "v-header",
    components: {
    },
    data() {
        return {
        }
    },
    methods: {
        logout () {
            console.log('logout')
            this.$http.get('/session/logout')
        }
    }
};
</script>

<style lang="css" scoped>
    @import "../assets/css/Header.scss";
</style>
