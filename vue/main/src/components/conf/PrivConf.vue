<template>
    <div>
        <div class="filter">
            <el-input class="input" v-on:input="onInput" v-model.number="phone" placeholder="请输入用户手机号码" @keyup.enter.native="query()" clearable></el-input>
            <el-button class="button" type="primary" @click="query()" accesskey="1">查询</el-button>
            <el-button class="button" type="primary" @click="submit()" accesskey="1">提交</el-button>
        </div>
        <div class="select">
            <el-radio-group class="radiogroup" v-model.number="admin">
                <el-radio :label="0" type="number">普通用户</el-radio>
                <el-radio :label="1" type="number">超级管理员</el-radio>
            </el-radio-group>
        </div>
        <div>
            <el-tree
                :data="all"
                show-checkbox
                default-expand-all
                node-key="id"
                ref="tree"
                :props="defaultProps">
            </el-tree>
        </div>
        </el-col>
    </div>
</template>

<script>
    import { mapState } from 'vuex'

    export default {
        name: "v-select",
        data() {
            return {
                admin: 0,
                phone: '',
                userId: '',
                selected: {},
                editdisable: true,
                defaultProps: {
                    children: 'children',
                    label: 'label'
                }
            }
        },
        props: {
        },
        methods: {
            onInput(val) {
                this.userId = 0
            },
            query() {
                this.$http.get('/api/user/' + this.phone)
                    .then(response => {
                        if (parseInt(response.data.status, 10) == 0) {
                            this.userId = response.data.data.user_id
                            let checked = response.data.data.access_level.split(";")
                            console.log(checked)
                            this.$refs.tree.setCheckedKeys(checked);
                            console.log(this.userId)
                        }
                    })
            },
            EnableEdit() {
                this.editdisable = false
            },
            DisableEdit() {
                this.editdisable = true
            },
            submit() {
                this.DisableEdit()
                console.log(this.$refs.tree.getCheckedNodes());
                let accessLevel = ''
                if (this.admin == 1) {
                    for (let i = 0; i < this.items.length; i++) {
                        let data = this.items[i]
                        accessLevel += data['id'] + ';'
                    }
                }
                else {
                    for (let i = 0; i < this.$refs.tree.getCheckedNodes().length; i++) {
                        let data = this.$refs.tree.getCheckedNodes()[i]
                        accessLevel += data['id'] + ';'
                    }
                }
                this.$http.post('/api/grant', {'user_id':this.userId,'access_level':accessLevel})
                console.log(this.all)
            }
        },
        computed: mapState({
            all: (state) => {
                    console.log(state.itemstree)
                    for (let i in state.itemstree) {
                        let item = state.itemstree[i]
                        item['label'] = item['name']

                        for (let j in item.items) {
                            let item1 = item.items[j]
                            item1['label'] = item1['name']
                        }
                        item['children'] = item['items']
                    }
                    return state.itemstree
                },
            checkd: (state) => {
                    console.log(state.itemstree)
                    for (let i in state.itemstree) {
                        let item = state.itemstree[i]
                        item['label'] = item['name']

                        for (let j in item.items) {
                            let item1 = item.items[j]
                            item1['label'] = item1['name']
                        }
                        item['children'] = item['items']
                    }
                    return state.itemstree
                },
                items:"items",
                loading:"loading"
            })
    }
</script>

<style lang="scss" scoped>
.filter {
    margin: 5px;
    border-radius: 4px;
    font-size: 14px;
    .input {
        width: 200px;
    }
    .button {
        width: 80px;
    }
}
.select {
    margin-top: 5px;
    margin-bottom: 20px;
    border-radius: 4px;
}
.input {
    border-radius: 4px;
    padding: 10px 0;
    background-color: #f9fafc;
    margin-top: 10px;
    margin-bottom: 10px;
}
</style>

