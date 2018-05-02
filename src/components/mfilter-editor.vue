<template lang="pug">
  v-container
    v-layout(row, wrap)
      v-flex(xs12)
        v-tabs
          v-tab(
            @click="handleTabClick(0)"
            :disabled="parseError == '' ? false : true"
          ) {{$t('m_filter_editor.tabs.visual_editor')}}
          v-tab(
            @click="handleTabClick(1)"
          ) {{$t('m_filter_editor.tabs.advanced_editor')}}

    template(v-if="activeTab === 0")
      filter-group(
        initialGroup,
        :index = 0,
        :condition.sync="filter[0].condition",
        :possibleFields="possibleFields",
        :rules="filter[0].rules",
        :groups="filter[0].groups",
      )

    template(v-if="activeTab === 1")
      v-text-field(
        ref="input",
        v-model="inputValue",
        rows="20",
        :label="$t('m_filter_editor.tabs.advanced_editor')",
        textarea
      )
      v-btn(@click="handleParseClick") {{$t('common.parse')}}
      p(v-if="parseError !== ''") {{ parseError }}
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import FilterGroup from '@/components/other/mfilter-editor/filter-group.vue';

const { mapGetters, mapActions } = createNamespacedHelpers('MFilterEditor');

export default {
  name: 'mfilter-editor',

  components: {
    FilterGroup,
  },

  data() {
    return {
      newRequest: '',
    };
  },

  computed: {
    ...mapGetters(['filter2request', 'filter', 'possibleFields', 'activeTab', 'parseError']),

    /**
     * @description Value of the input field of the advanced editor.
     * Prettify the value of the parsed filter ('filter2request')
     */
    inputValue: {
      get() {
        return JSON.stringify(this.filter2request, undefined, 4);
      },
      set(newVal) {
        this.newRequest = newVal;
      },
    },
  },

  methods: {
    ...mapActions(['changeActiveTab', 'updateFilter', 'onParseError', 'deleteParseError']),

    handleTabClick(tab) {
      this.newRequest = '';
      this.changeActiveTab(tab);
    },

    handleParseClick() {
      this.deleteParseError();
      try {
        if (this.newRequest === '') {
          this.updateFilter(JSON.parse(JSON.stringify(this.filter2request)));
          return this;
        }
        this.updateFilter(JSON.parse(this.newRequest));
        return this;
      } catch (e) {
        this.onParseError(e.message);
        return e;
      }
    },
  },
};
</script>
