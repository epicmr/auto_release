<template>
    <div>
        <el-checkbox-group v-model="checkedApp">
            <el-row v-if="list_type === parseInt(app.serv_type)" v-for="app in this.$store.state.data">
                <el-col :span=12>
                    <el-checkbox name="appn" :label="app.serv_name" :key="app.serv_name" style="margin-top: 10px" size="medium"></el-checkbox>
                </el-col>
                <el-col :span=12>
                    <el-row v-for="time in app.serv_state">
                        <el-col :span=12>
                            <small>{{time.host_name}}</small>
                        </el-col>
                        <el-col :span=12>
                            <small>{{time.serv_time}}</small>
                        </el-col>
                    </el-row>
                </el-col>
            </el-row>
        </el-checkbox-group>
    </div>
</template>

<script>
    export default {
        name: "v-servlist",
        props : {
            list_type : {
                type : Number,
                default : 0
            }
        },
        data() {
            return {
                parent:this,
                checkedApp:[]
            }
        }
    }
</script>
<style scoped>
small{
    font-size: 12px;
    color: #66b1ff;
    vertical-align: bottom;
    display: inline-block;
    margin-top: 10px;
}
</style>

