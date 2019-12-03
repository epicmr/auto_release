<template>
    <div class="home">
        <router-view/>
    </div>
</template>

<script>

export default {
    data() {
        return {
            parent: this
        }
    },
    created: function() {
        //console.log(this.$store.state.loadingcount)
    },
    name: "v-welcome"
};
</script>

<style lang="scss">
.home{
    position: absolute;
    left: 50%;
    top: 50%;
}
<style>
