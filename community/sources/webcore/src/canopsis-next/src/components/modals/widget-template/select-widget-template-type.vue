<template lang="pug">
  modal-wrapper(close)
    template(#title="")
      span {{ $t('modals.selectWidgetTemplateType.title') }}
    template(#text="")
      v-layout(column)
        v-card.my-1.cursor-pointer(
          v-for="{ value, text, icon } in availableTypes",
          :key="value",
          @click="selectType(value)"
        )
          v-card-title(primary-title)
            v-layout(wrap, justify-between)
              v-flex(xs11)
                div.subheading {{ text }}
              v-flex
                v-icon {{ icon }}
</template>

<script>
import { MODALS, WIDGET_TEMPLATES_TYPES, COLUMNS_WIDGET_TEMPLATES_TYPES } from '@/constants';

import { modalInnerMixin } from '@/mixins/modal/inner';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Modal to create widget
 */
export default {
  name: MODALS.createWidgetTemplate,
  components: { ModalWrapper },
  mixins: [modalInnerMixin],
  computed: {
    availableTypes() {
      return Object.values(WIDGET_TEMPLATES_TYPES).map(type => ({
        value: type,
        icon: COLUMNS_WIDGET_TEMPLATES_TYPES.includes(type) ? 'view_week' : 'description',
        text: this.$t(`widgetTemplate.types.${type}`),
      }));
    },
  },
  methods: {
    selectType(type) {
      this.$modals.show({
        name: MODALS.createWidgetTemplate,
        config: {
          title: this.$t('modals.createWidgetTemplate.create.title'),
          widgetTemplate: { type },
          action: this.config.action,
        },
      });

      this.$modals.hide();
    },
  },
};
</script>
