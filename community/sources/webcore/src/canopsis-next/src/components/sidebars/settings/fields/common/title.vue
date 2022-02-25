<template lang="pug">
  v-list-group
    template(#activator="")
      v-list-tile {{ title }}
        div.font-italic.caption.ml-1 ({{ $t('common.optional') }})
    v-container
      v-text-field(
        v-field="value",
        :placeholder="$t('settings.widgetTitle')"
      )
</template>

<script>
/**
 * Component to set the title of a widget
 *
 * @prop {String} [value] - Value of the title
 */
export default {
  props: {
    value: {
      type: String,
      required: false,
    },
    title: {
      type: String,
      required: false,
    },
  },
};
</script>
