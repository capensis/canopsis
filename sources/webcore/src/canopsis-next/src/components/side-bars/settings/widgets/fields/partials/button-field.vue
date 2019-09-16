<template lang="pug" functional>
  v-container.pa-3(fluid)
    v-layout(align-center, justify-space-between)
      slot(name="title")
      v-layout(justify-end)
        v-btn.primary(data-test="createButton", v-if="props.isEmpty", small, @click="listeners.create")
          span {{ parent.$t('common.create') }}
        template(v-else)
          v-btn.primary(data-test="editButton", small, @click="listeners.edit")
            span {{ parent.$t('common.edit') }}
          v-btn.error(small, @click="listeners.delete")
            v-icon(data-test="deleteButton") delete
</template>

<script>
export default {
  props: {
    isEmpty: {
      type: Boolean,
      default: true,
    },
  },
};
</script>

