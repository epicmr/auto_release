<template>
    <div>
        <el-row>
            <el-button type="primary" @click="Pack()" plain>打包</el-button>
            <el-button type="primary" @click="Trans()" plain>上传</el-button>
            <el-button type="primary" @click="Post()" plain>发布</el-button>
        </el-row>
        <el-row>
            <el-button type="primary" @click="OneKeyPost()" plain>一键发布</el-button>
        </el-row>
    </div>
</template>

<script>
    export default {
        name: "v-button",
        data() {
            return {
                parent:this
            }
        },
        methods : {
            Pack() {
                console.log("pack serv")
                let promises = [];
                let app = document.forms[0].serv_list;
                for (var i = 0; i < app.length; ++i) {
                    if (app[i].checked){
                        let serv_name = app[i].value
                        console.log("pack serv", serv_name)
                        promises.push(this.$store.commit("Pack", {"serv_name":serv_name,"env":this.$store.state.menuvalue}))
                    }
                };
                return Promise.all(promises)
            },
            Trans() {
                console.log("trans serv")
                let promises = [];
                var app = document.forms[0].serv_list;
                for (var i = 0; i < app.length; ++i) {
                    if (app[i].checked){
                        let serv_name = app[i].value
                        console.log("trans serv", serv_name)
                        promises.push(this.$store.commit("Trans", {"serv_name":serv_name,"env":this.$store.state.menuvalue}))
                    }
                };
                return Promise.all(promises)
            },
            Post() {
                console.log("post serv")
                let promises = [];
                var app = document.forms[0].serv_list;
                for (var i = 0; i < app.length; ++i) {
                    if (app[i].checked){
                        let serv_name = app[i].value
                        console.log("post serv", serv_name)
                        promises.push(this.$store.commit("Post", {"serv_name":serv_name,"env":this.$store.state.menuvalue}))
                    }
                }
                return Promise.all(promises)
            },
            async OneKeyPost() {
                await this.Pack()
                await this.Trans()
                await this.Post()
            }
        }
    };
</script>

<style scoped>
.el-row {
    margin-bottom: 20px;
}
</style>
