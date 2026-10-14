<template>
  <v-row class="theme-form__colors" no-gutters>
    <v-col cols="12">
      <theme-color-font-size-field v-model="fontSizeField" />
    </v-col>
  </v-row>
</template>

<script>
import { useFormModels } from '@/hooks/form';

import ThemeColorFontSizeField from './fields/theme-color-font-size-field.vue';

export default {
  components: { ThemeColorFontSizeField },

  setup() {
    const { fontSizeField } = useFormModels({ fontSizeField: 'font_size' });

    return { fontSizeField };
  },
};
</script>
