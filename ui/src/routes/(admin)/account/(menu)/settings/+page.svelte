<script lang="ts">
	import { getContext } from 'svelte';
	import type { Writable } from 'svelte/store';
	import SettingsModule from './settings_module.svelte';

	let adminSection: Writable<string> = getContext('adminSection');
	adminSection.set('settings');

	interface Props {
		profile: {
			full_name: string;
			company_name: string;
			website: string;
			unsubscribed: string;
		};
		user: {
			email: string;
		};
	}

	let { profile, user }: Props = $props();
</script>

<svelte:head>
	<title>Settings</title>
</svelte:head>

<h1 class="text-2xl font-bold mb-6">Settings</h1>

<SettingsModule
	title="Profile"
	editable={false}
	fields={[
		{ id: 'fullName', label: 'Name', initialValue: profile?.full_name ?? '' },
		{
			id: 'companyName',
			label: 'Company Name',
			initialValue: profile?.company_name ?? ''
		},
		{
			id: 'website',
			label: 'Company Website',
			initialValue: profile?.website ?? ''
		}
	]}
	editButtonTitle="Edit Profile"
	editLink="/account/settings/edit_profile"
/>

<SettingsModule
	title="Email"
	editable={false}
	fields={[{ id: 'email', initialValue: user?.email || '' }]}
	editButtonTitle="Change Email"
	editLink="/account/settings/change_email"
/>

<SettingsModule
	title="Security"
	editable={false}
	fields={[{ id: 'password', initialValue: '••••••••••••••••' }]}
	editButtonTitle="Change Security Settings"
	editLink="/account/settings/security"
/>

<SettingsModule
	title="Email Subscription"
	editable={false}
	fields={[
		{
			id: 'subscriptionStatus',
			initialValue: profile?.unsubscribed ? 'Unsubscribed' : 'Subscribed'
		}
	]}
	editButtonTitle="Change Subscription"
	editLink="/account/settings/change_email_subscription"
/>

<SettingsModule
	title="Danger Zone"
	editable={false}
	dangerous={true}
	fields={[]}
	editButtonTitle="Delete Account"
	editLink="/account/settings/delete_account"
/>
