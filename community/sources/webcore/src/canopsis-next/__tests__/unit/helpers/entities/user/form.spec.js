import { describe, expect, it } from 'vitest';

import { formToUserRequest, userToForm } from '@/helpers/entities/user/form';

describe('user form helpers', () => {
  it('does not expose the system-managed ID in the form', () => {
    const form = userToForm({ _id: 'user-id' });

    expect(form).not.toHaveProperty('_id');
  });

  it('omits the system-managed ID from the request', () => {
    const user = formToUserRequest({
      _id: 'user-id',
      name: 'username',
      roles: [],
    });

    expect(user).toEqual({ name: 'username', roles: [] });
  });
});
