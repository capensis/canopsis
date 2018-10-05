<template lang="pug">
  v-card
    v-card-title.blue.darken-4.white--text
      h2 Manage stats groups
    v-layout
      v-flex(xs4)
        v-card.my-1
          v-card-title
            h2 {{ editing ? $t('common.edit') : $t('common.add') }}
            v-btn(v-if="editing", @click="editing = false")
              v-icon close
            v-form
              v-text-field(:placeholder="$t('common.title')", v-model="form.title")
              v-layout(wrap, column)
                v-flex(xs6)
                  v-btn(@click="showFilterModal", small) {{ $t('settings.filterEditor') }}
              v-btn(@click="addGroup").green.darken-4.white--text.mt-3 {{ $t('common.save') }}
      v-flex(xs8)
        v-container.pt-0
          v-card.my-1(v-for="(group, index) in groups", :key="index")
            v-layout(align-center, justify-between)
              v-flex
                div.ml-2 {{ group.title }}
              v-flex(xs4)
                v-layout
                  v-btn(@click="editGroup(group, index)", fab, small, depressed)
                    v-icon edit
                  v-btn(@click="deleteGroup(index)", fab, small, depressed)
                    v-icon delete
    v-layout(justify-end)
      v-btn(@click="save").green.darken-4.white--text.mt-3 {{ $t('common.save') }}
</template>

<script>
import modalInnerMixin from '@/mixins/modal/modal-inner';
import { MODALS } from '@/constants';

export default {
  name: MODALS.manageHistogramGroups,
  mixins: [modalInnerMixin],
  data() {
    return {
      editing: false,
      editingGroupIndex: null,
      form: {
        title: '',
        filter: '',
      },
      groups: [],
    };
  },
  mounted() {
    this.groups = [...this.config.groups];
  },
  methods: {
    showFilterModal() {
      this.showModal({
        name: MODALS.createFilter,
        config: {
          title: 'modals.filter.create.title',
          filter: this.form.filter,
          action: filter => this.form.filter = filter,
        },
      });
    },
    editGroup(group, index) {
      this.editing = true;
      this.editingGroupIndex = index;
      this.form = { ...group };
    },
    deleteGroup(index) {
      this.$delete(this.groups, index);
    },
    addGroup() {
      if (this.editing) {
        // Using Vue.set to be sure the groups list will be updated + provoke a re-render of the list
        this.$set(this.groups, this.editingGroupIndex, { ...this.form });
        this.editing = false;
      } else {
        this.groups.push({ ...this.form });
      }
    },
    async save() {
      await this.config.action(this.groups);
      this.hideModal();
    },
  },
};
</script>
