import { redirect } from '@sveltejs/kit';

export const load = async ({ locals }) => {

  if (!locals.pb.authStore.isValid) {
    throw redirect(303, '/login');
  }

  // GET USERS
  const users = await locals.pb.collection('users').getFullList({
    sort: '-created'
  });

  console.log({ users })

  return { users: users };
};