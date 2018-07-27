<template lang="pug">
  div.attribute
    div.attribute-name.text-md-right
      b
        slot(name="name")
    div.attribute-content
      slot(name="content")
</template>

<style scoped>
  .attribute-name {
    width: 30%;
    display: inline-block;
  }

  .attribute-content {
    display: inline-block;
  }

  .attribute-content {
    padding-left: 10px;
  }
</style>
