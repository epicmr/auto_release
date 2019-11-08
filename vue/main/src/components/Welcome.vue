<template>
    <div>
        <h1 class="center">Welcome</h1>
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
        console.log(this.$store.state.loadingcount)
    },
    name: "v-welcome"
};
</script>

<style>
.center{
    text-align:center;
}
<style>
