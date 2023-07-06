import { TAG_TYPES } from '@/constants';

export default {
  importedTag: 'Tag is imported',
  types: {
    [TAG_TYPES.imported]: 'Imported',
    [TAG_TYPES.created]: 'Created',
  },
  deleteConfirmation: 'Are you sure you want to delete this tag?\nThis action cannot be canceled. The tag will be removed from all alarms affected.'
    + ' | Are you sure you want to delete this tags?\nThis action cannot be canceled. The tags will be removed from all alarms affected.',
};
