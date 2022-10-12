<template lang="pug">
  widget-settings-item
    template(#title="")
      span {{ title }}
      v-icon.ml-2 {{ sizeIcon }}
    c-column-size-field(v-field="value", :mobile="mobile", :tablet="tablet")
</template>

<script>
import WidgetSettingsItem from '@/components/sidebars/settings/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem },
  props: {
    value: {
      type: Number,
      required: false,
    },
    title: {
      type: String,
      required: true,
    },
    mobile: {
      type: Boolean,
      default: false,
    },
    tablet: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    sizeIcon() {
      if (this.mobile) {
        return 'phone_android';
      }

      if (this.tablet) {
        return 'tablet';
      }

      return 'laptop';
    },
  },
};
</script>
