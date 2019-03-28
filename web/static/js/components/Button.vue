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
                console.log(document.forms[0])
                console.log(document.forms[0].appn)
                let app = document.forms[0].appn;
                console.log(app.length)
                for (var i = 0; i < app.length; ++i) {
                    console.log(app[i].checked)
                    if (app[i].checked){
                        var serv_name = this.$store.state.data[app[i].value]["serv_name"]
                        console.log("pack serv", serv_name)
                        promises.push(this.$store.dispatch("Pack", {"serv_name":serv_name,"env":this.$store.state.menuvalue}))
                    }
                };
                return Promise.all(promises)
            },
            Trans() {
                let promises = [];
                var app = document.forms[0].appn;
                for (var i = 0; i < app.length; ++i) {
                    if (app[i].checked){
                        var serv_name = this.$store.state.data[app[i].value]["serv_name"]
                        console.log("trans serv", serv_name)
                        promises.push(this.$store.dispatch("Trans", {"serv_name":serv_name,"env":this.$store.state.menuvalue}))
                    }
                };
                return Promise.all(promises)
            },
            Post() {
                let promises = [];
                var app = document.forms[0].appn;
                for (var i = 0; i < app.length; ++i) {
                    if (app[i].checked){
                        var serv_name = this.$store.state.data[app[i].value]["serv_name"]
                        console.log("post serv", serv_name)
                        promises.push(this.$store.dispatch("Post", {"serv_name":serv_name,"env":this.$store.state.menuvalue}))
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
