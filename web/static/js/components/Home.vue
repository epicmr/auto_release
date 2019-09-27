<template>
    <div id="home" class="home">
        <el-row>
            <v-header></v-header>
        </el-row>
        <el-row>
            <el-col :span="3" :offset="4" class="site-sidemenu">
                <v-sidemenu></v-sidemenu>
            </el-col>
            <el-col :span="12" :offset="1" class="site-main" element-loading-text="拼命发布中" v-loading="this.$store.state.loading">
                <router-view></router-view>
            </el-col>
        </el-row>
    </div>
</template>

<script>
import VSidemenu from "./Sidemenu";
import VLoading from "./Loading";
import VHeader from "./Header";

export default {
    components: {
        VHeader,
        VLoading,
        VSidemenu
    },
    data() {
        return {
            checkedApp:[]
        }
    },
    methods: {
    },
    name: "home"
};
</script>

<style>
.home {
    width: 100%;
    height: 100%;
}
.grid-content {
    border-radius: 4px;
    min-height: 36px;
    margin-top: 30px;
}
.row-bg {
    padding: 10px 0;
    background-color: #f9fafc;
}
.line-none {
    text-decoration: none;
}
.site-sidemenu {
    min-height: calc(100vh - 5rem);
    padding-top: 50px;
    padding-bottom: 50px;
    padding-left: 0;
}
.site-main {
    min-height: calc(100vh - 5rem);
    padding-top: 50px;
    padding-bottom: 50px;
    padding-right: 0;
}
</style>

