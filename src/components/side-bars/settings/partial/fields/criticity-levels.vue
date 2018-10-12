<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{$t('settings.criticityLevels')}}
    v-container
      v-layout(wrap)
        v-flex(xs12)
          v-layout(align-center, justify-space-around)
            div Minor :
            v-flex(xs3)
              v-text-field(
                type="number",
                :value="value.minor",
                @input="updateField('minor', parseInt($event, 10))"
                )
        v-flex(xs12)
          v-layout(align-center, justify-space-around)
            div Major :
            v-flex(xs3)
              v-text-field(
              type="number",
              :value="value.major",
              @input="updateField('major', parseInt($event, 10))"
              )
        v-flex(xs12)
          v-layout(align-center, justify-space-around)
            div Critical :
            v-flex(xs3)
              v-text-field(
              type="number",
              :value="value.critical",
              @input="updateField('critical', parseInt($event, 10))"
              )
</template>

<script>
import formMixin from '@/mixins/form';

export default {
  mixins: [formMixin],
  props: {
    value: {
      type: Object,
    },
  },
};
</script>
