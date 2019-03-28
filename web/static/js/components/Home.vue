<template>
    <div id="home" class="grid-content">
        <el-col :span="3" :offset="3">
            <v-sidemenu></v-sidemenu>
        </el-col>
        <el-col :span="16" :offset="2" element-loading-text="拼命发布中" v-loading="this.$store.state.loading">
            <router-view></router-view>
        </el-col>
    </div>
</template>

<script>
import VSidemenu from "./Sidemenu";
import VLoading from "./Loading";

export default {
    components: {
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
.el-row {
    margin-bottom: 1px;
}
.el-col {
    border-radius: 4px;
}
.bg-purple-dark {
    background: #99a9bf;
}
.bg-purple {
    background: #d3dce6;
}
.bg-purple-light {
    background: #e5e9f2;
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
</style>

