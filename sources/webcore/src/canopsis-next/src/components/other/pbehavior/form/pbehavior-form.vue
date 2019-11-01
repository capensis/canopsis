<template lang="pug">
  div
    pbehavior-general-form(v-field="form.general", :noFilter="noFilter")
    pbehavior-exdates-form.mt-2(v-show="form.general.rrule", v-field="form.exdate")
    pbehavior-comments-form.mt-2(v-field="form.comments", :author="author")
</template>

<script>
import formMixin from '@/mixins/form';

import PbehaviorGeneralForm from './partials/pbehavior-general-form.vue';
import PbehaviorExdatesForm from './partials/pbehavior-exdates-form.vue';
import PbehaviorCommentsForm from './partials/pbehavior-comments-form.vue';

export default {
  components: {
    PbehaviorGeneralForm,
    PbehaviorExdatesForm,
    PbehaviorCommentsForm,
  },
  mixins: [formMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
    author: {
      type: String,
    },
    noFilter: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
