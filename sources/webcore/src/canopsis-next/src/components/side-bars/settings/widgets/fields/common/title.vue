<template lang="pug">
  v-list-group(data-test="settingWidgetTitle")
    v-list-tile(slot="activator") {{ title }}
      .font-italic.caption.ml-1 ({{ $t('common.optional') }})
    v-container
      v-text-field(
        data-test="widgetTitleField",
        :value="value",
        @input="$emit('input', $event)",
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
    value: String,
    title: String,
  },
};
</script>
