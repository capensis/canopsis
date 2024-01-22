<template>
  <c-advanced-data-table
    :headers="headers"
    :items="users"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    :select-all="removable"
    advanced-pagination
    search
    @update:options="$emit('update:options', $event)"
  >
    <template #mass-actions="{ selected }">
      <c-action-btn
        v-if="removable"
        type="delete"
        @click="$emit('remove-selected', selected)"
      />
    </template>
    <template #name="{ item }">
      {{ item.display_name }}
    </template>
    <template #enable="{ item }">
      <c-enabled :value="item.enable" />
    </template>
    <template #active="{ item }">
      <c-enabled :value="item.active_connects > 0" />
    </template>
    <template #source="{ item }">
      {{ item.source || $constants.AUTH_METHODS.local }}
    </template>
    <template #roles="{ item }">
      <v-chip-group>
        <v-chip
          v-for="role in item.roles"
          :key="role._id"
        >
          {{ role.name }}
        </v-chip>
      </v-chip-group>
    </template>
    <template #actions="{ item }">
      <v-layout>
        <c-action-btn
          v-if="updatable"
          type="edit"
          @click.stop="$emit('edit', item)"
        />
        <c-action-btn
          v-if="removable"
          type="delete"
          @click.stop="$emit('remove', item)"
        />
      </v-layout>
    </template>
  </c-advanced-data-table>
</template>

<script>
export default {
  props: {
    users: {
      type: Array,
      default: () => [],
    },
    totalItems: {
      type: Number,
      required: false,
    },
    pending: {
      type: Boolean,
      default: false,
    },
    options: {
      type: Object,
      required: true,
    },
    removable: {
      type: Boolean,
      default: false,
    },
    updatable: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    headers() {
      return [
        {
          text: this.$t('common.username'),
          value: 'name',
        },
        {
          text: this.$t('user.firstName'),
          value: 'firstname',
          sortable: false,
        },
        {
          text: this.$t('user.lastName'),
          value: 'lastname',
          sortable: false,
        },
        {
          text: this.$tc('common.role', 2),
          value: 'roles',
          sortable: false,
        },
        {
          text: this.$t('user.active'),
          value: 'active',
          sortable: false,
        },
        {
          text: this.$t('common.enabled'),
          value: 'enable',
        },
        {
          text: this.$t('user.auth'),
          value: 'source',
        },
        {
          text: this.$t('user.activeConnects'),
          value: 'active_connects',
          sortable: false,
        },
        {
          text: this.$t('common.actionsLabel'),
          value: 'actions',
          sortable: false,
        },
      ];
    },
  },
};
</script>
