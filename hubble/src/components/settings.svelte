<script>
    export let netname
    let newgroup;
    let form = "readonly";
    export let network;
    import { navigate } from "svelte-routing";
    import { createEventDispatcher } from 'svelte';
    const dispatch = createEventDispatcher();
    import API from '../api.js';
    import { getNotificationsContext } from 'svelte-notifications';
    const { addNotification } = getNotificationsContext();
    export const updateNetwork = async() => {
        try {
            const response = await API.post(`/networks/${network.name}/update`,
                {
                    'name': network.name,
                    'cidr': network.cidr,
                    'cipher': network.cipher,
                    'groups': network.groups,
                });
            console.log(response)
            addNotification({
                text: `Updated network ${network.name}`,
                position: 'top-right',
                type: 'success',
                removeAfter: 2000,
            })
            return response;
        } catch (error) {
            console.error(error);
            addNotification({
                text: `Could not update network ${network.name}`,
                position: 'top-right',
                type: 'danger',
                removeAfter: 2000,
            })
        }
    };
    async function handleClick(){
        await updateNetwork();
        dispatch('updateNetworkInfo');
        console.log("FORM SUBMITTED", network.cipher, network.cidr, network.name);
    }
    export const deleteNetwork = async() => {
        try {
            const response = await API.delete(`/networks/${netname}/delete`);
            console.log(response)
            dispatch('updateNetworks');
            navigate('/');
            addNotification({
                text: `Deleted network ${network.name}`,
                position: 'top-right',
                type: 'success',
                removeAfter: 2000,
            })
            return response;
        } catch (error) {
            console.error(error);
            addNotification({
                text: `Could not delete network ${network.name}`,
                position: 'top-right',
                type: 'danger',
                removeAfter: 2000,
            })
        }
    };
    $: {
        if (network.name != "") {
            form = "" // enable form
        }
    }
    function addGroup(){
        console.log("Adding", newgroup, "to", network.groups)
        var groups = network.groups
        groups.push(newgroup);
        network.groups = groups;
        newgroup = "";
    }
    function removeGroup(groupname){
        console.log("Removing", groupname, "from", network.groups);
        var groups = network.groups;
        for( var i = 0; i < groups.length; i++){ 
            if ( groups[i] === groupname) { 
                groups.splice(i, 1); 
            }
        }
        network.groups = groups;
    }
</script>
<main class="shadow p-4 mt-4">
    <h2 class="text-2xl">Network Settings</h2>
    <div id="cidr" class="flex flex-row my-4">
        <label for="netname" class="my-auto text-lg">Network CIDR: </label>
        <input type="text" id="name" name="name" bind:value={network.cidr} readonly={form}
            class="w-1/4 h-10 pl-3 pr-6 text-base border rounded-lg appearance-none focus:shadow-outline my-auto mx-4 bg-gray-100">
    </div>
    <div id="cipher" class="flex flex-row my-4">
        <label for="cipher" class="my-auto text-lg">Cipher Algorithm: </label>
        <select name="cipher" id="cipher" bind:value={network.cipher}
            class="w-1/4 h-10 pl-3 pr-6 text-base placeholder-gray-600 border rounded-lg appearance-none focus:shadow-outline my-auto mx-4">
            <option value="aes">AES</option>
            <option value="chachapoly">chacha</option>
        </select>
    </div>
    <div id="groups" class="bg-gray-200 rounded p-4">
        <h1 class="text-lg ml-4 mb-2">Groups</h1>
        <input type="text" id="newgroup" name="newgroup" bind:value={newgroup}
            class="w-1/4 h-10 pl-3 pr-6 text-base placeholder-gray-600 border rounded-lg appearance-none focus:shadow-outline my-auto ml-4 mr-2 bg-gray-100"
        />
        <button on:click={addGroup} class="bg-indigo-500 text-white rounded-xl p-1 px-2">+</button>
        <div class="flex flex-row flex-wrap pt-4 ml-2">
            {#each network.groups as group}
                <div class="bg-pink-400 rounded-xl text-white m-2 p-1 my-auto">
                    {group}
                    <button on:click={() => removeGroup(group)} class="outline-none border-none text-gray-500 mx-2 my-auto">x</button>
                </div>
            {/each}
        </div>
    </div>
    <div class="text-white my-4 flex flex-row-reverse">
        <button on:click|preventDefault={async() => {await deleteNetwork()}}
            class="bg-red-500 rounded-xl p-1 mx-4" {form}>
            Delete
        </button>
        <button on:click|preventDefault={async() => {await handleClick()}}
                class="bg-green-400 rounded-xl p-1 text-white" {form}>
            Update
        </button>
    </div>
</main>
