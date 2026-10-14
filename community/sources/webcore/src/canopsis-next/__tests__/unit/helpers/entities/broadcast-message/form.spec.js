import { BROADCAST_MESSAGE_VIEWS } from '@/constants';

import { viewsFormToMessage } from '@/helpers/entities/broadcast-message/form';

describe('viewsFormToMessage', () => {
  it('preserves an explicitly selected all-views parent', () => {
    const treeItems = {
      pages: [],
      views: [
        {
          value: BROADCAST_MESSAGE_VIEWS.allViews,
          children: [
            {
              value: 'view-group',
              children: [{ value: 'view' }],
            },
          ],
        },
      ],
      playlists: [],
    };

    expect(viewsFormToMessage({
      pages: [],
      views: [BROADCAST_MESSAGE_VIEWS.allViews],
      playlists: [],
    }, treeItems)).toContain(BROADCAST_MESSAGE_VIEWS.allViews);
  });
});
