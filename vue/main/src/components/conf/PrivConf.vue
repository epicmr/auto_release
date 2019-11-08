<template>
    <div class="grid-content" v-loading="this.loading">
        <el-col>
            <el-row>
                <el-col :span=16>
                    <el-input v-model="selected.name" class="input" :disabled="editdisable">
                        <template slot="prepend">用户</template>
                    </el-input>
                    <el-input v-model="selected.type" class="input" :disabled="editdisable">
                        <template slot="prepend">行为</template>
                    </el-input>
                    <el-input v-model="selected.group" class="input" :disabled="editdisable">
                        <template slot="prepend">资源</template>
                    </el-input>
                </el-col>
            </el-row>
            <el-row>
                <el-button class="button" type="primary" @click="EnableEdit" plain>{{AddOrEdit}}</el-button>
                <el-button class="button" type="primary" @click="Submit()" :disabled="editdisable" plain>提交</el-button>
            </el-row>
        </el-col>
    </div>
</template>

<script>
    import { mapState } from 'vuex'

    export default {
        name: "v-select",
        data() {
            return {
                selectedName: '',
                selected: {},
                editdisable: true
            }
        },
        props: {
        },
        methods: {
            EnableEdit() {
                this.editdisable = false
            },
            DisableEdit() {
                this.editdisable = true
            },
            Submit() {
                this.DisableEdit()
                let list = []
                list.push(this.selected)
                this.$store.commit("UserGroup", list)
            }
        },
        computed: {
            ...mapState([
                "loading"
            ]),
            AddOrEdit: {
                get() {
                    if (this.selected.id == 0 ) {
                        return "新增"
                    }
                    return "编辑"
                }
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

