<template>
    <div class="grid-content" v-loading="this.$store.state.loading">
            <el-row>
                <el-col :span=12>
                    <el-input v-model.number="userId" class="input" clearable>
                        <template slot="prepend">用户ID</template>
                    </el-input>
                <el-col>
            </el-row>
            <el-row>
                <el-button class="button" type="primary" @click="Conf()" plain>提交</el-button>
                <el-button class="button" type="primary" @click="Conf()" plain>查询</el-button>
            </el-row>
        </el-col>
    </div>
</template>

<script>
    export default {
        name: "v-userrepair",
        data() {
            return {
                userId : ''
            }
        },
        methods: {
            Conf() {
                console.log(this.userId)
                this.$store.commit("RepairUser", this.userId)
            }
        }
    }
</script>

<style scoped>
.input {
    border-radius: 4px;
    padding: 10px 0;
    background-color: #f9fafc;
    margin-top: 10px;
    margin-bottom: 10px;
}
</style>

