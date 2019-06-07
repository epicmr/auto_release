<template>
    <div style="margin-top: 20px">
        <el-row>
            <el-button type="primary" @click="Pack()" plain>打包</el-button>
            <el-button type="primary" @click="Trans()" plain>上传</el-button>
            <el-button type="primary" @click="Post()" plain>发布</el-button>
        </el-row>
        <el-row>
            <el-button type="primary" @click="OneKeyPost()" plain>一键发布</el-button>
        </el-row>
        <el-row>
            <el-button type="primary" @click="CheckMD5()" plain>检查发布MD5</el-button>
            <el-button type="primary" @click="CheckTime()" plain>检查发布时间</el-button>
        </el-row>
    </div>
</template>

<script>
    import { mapState } from 'vuex'
    var async = require('async');
    require('babel-polyfill');

    export default {
        name: "v-button",
        data() {
            return {
                parent:this
            }
        },
        computed: {
            ...mapState([
                "menuvalue"
            ])
        },
        methods : {
            Request(uri, data) {
                return this.$http.post(uri, data)
            },
            open(msg) {
                console.log(msg)
                if (msg.length > 0) {
                    msg = msg.replace(/\[/g, "<b>[")
                    msg = msg.replace(/\]/g, "]</b>")
                    this.$msgbox.alert(msg, '检测到发布错误', {
                        dangerouslyUseHTMLString: true,
                        confirmButtonText: '确定',
                    });
                }
            },
            CheckMD5() {
                let promises = [];
                let app = document.forms[0].serv_list;
                for (var i = 0; i < app.length; ++i) {
                    if (app[i].checked){
                        let serv_name = app[i].value
                        promises.push(this.Request('/api/checkmd5', {"serv_name":serv_name,"env":this.menuvalue}))
                    }
                };
                Promise.all(promises)
                    .then(response => {
                        console.log(response)
                        let msg = ''
                        for (let i in response) {
                            let data = response[i].data
                            console.log(data.status)
                            if (data.status != 0) {
                                msg += "<div>" + data.message + "</div>"
                            }
                        }
                        this.open(msg)
                    })
            },
            CheckTime() {
                let promises = [];
                let app = document.forms[0].serv_list;
                for (var i = 0; i < app.length; ++i) {
                    if (app[i].checked){
                        let serv_name = app[i].value
                        promises.push(this.Request('/api/checktime', {"serv_name":serv_name,"env":this.menuvalue}))
                    }
                };
                Promise.all(promises)
                    .then(response => {
                        console.log(response)
                        let msg = ''
                        for (let i in response) {
                            let data = response[i].data
                            console.log(data.status)
                            if (data.status != 0) {
                                msg += "<div>" + data.message + "</div>"
                            }
                        }
                        this.open(msg)
                    })
            },
            Pack() {
                console.log("pack serv")
                let promises = [];
                let app = document.forms[0].serv_list;
                for (var i = 0; i < app.length; ++i) {
                    if (app[i].checked){
                        let serv_name = app[i].value
                        console.log("pack serv", serv_name)
                        promises.push(this.Request('/release/pack', {"serv_name":serv_name,"env":this.menuvalue}))
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
                        promises.push(this.Request('/release/trans', {"serv_name":serv_name,"env":this.menuvalue}))
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
                        promises.push(this.Request('/release/post', {"serv_name":serv_name,"env":this.menuvalue}))
                    }
                }
                return Promise.all(promises)
            },
            async OneKeyPost() {
                this.$store.commit("ShowLoading")
                await this.Pack()
                await this.Trans()
                await this.Post()
                this.$store.commit("HiddenLoading")
            }
        }
    };
</script>

<style scoped>
.el-row {
    margin-bottom: 20px;
}
</style>
