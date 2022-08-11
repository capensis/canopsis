<template lang="pug">
  v-layout(column)
    c-number-field(
      v-field="value.lat",
      :label="$t('map.latitude')",
      :name="`${name}.latitude`",
      :disabled="disabled",
      :min="-90",
      :max="90",
      :step="step",
      required,
      @paste="pasteHandler('lat', $event)"
    )
    c-number-field(
      v-field="value.lng",
      :label="$t('map.longitude')",
      :name="`${name}.longitude`",
      :disabled="disabled",
      :min="-180",
      :max="180",
      :step="step",
      required,
      @paste="pasteHandler('lng', $event)"
    )
</template>

<script>
import { isNumber, isString } from 'lodash';

import { formMixin } from '@/mixins/form';

export default {
  mixins: [formMixin],
  props: {
    value: {
      type: Object,
      required: false,
    },
    name: {
      type: String,
      default: 'coordinate',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    step: {
      type: Number,
      default: 0.01,
    },
  },
  methods: {
    pasteHandler(name, event) {
      const data = (event.clipboardData || window.clipboardData).getData('text');

      if (data && isString(data)) {
        if (data.includes(',')) {
          const [latString, lngString] = data.split(',');
          const lat = +latString;
          const lng = +lngString;

          if (isNumber(lat) && isNumber(lng)) {
            event.preventDefault();

            this.updateModel({ lat, lng });
          }
        }
      }
    },
  },
};
</script>
