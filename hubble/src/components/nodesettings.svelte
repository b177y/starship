<script>
    export let netname
    export let nodename
    export let netgroups
    export let netgroupOptions = netgroups
    let newgroup
    let node = {
        is_lighthouse: false,
        name: '',
        hostname: '',
        address: '',
        static_address: '',
        listen_port: 0,
        firewall_outbound: [],
        firewall_inbound: [],
        groups: [],
    }
    import { getNotificationsContext } from 'svelte-notifications';
    const { addNotification } = getNotificationsContext();
    import API from '../api.js';
    import FirewallRule from './firewall.svelte';
    export const updateNode = async() => {
        try {
            const response = await API.post(`/networks/${netname}/nodes/${nodename}/update`,
                {
                    'is_lighthouse': node.is_lighthouse,
                    'static_address': node.static_address,
                    'listen_port': node.listen_port,
                    'groups': node.groups,
                    'firewall_inbound': node.firewall_inbound,
                    'firewall_outbound': node.firewall_outbound,
                });
            console.log(response)
            addNotification({
                text: `Updated node ${node.name}`,
                position: 'top-right',
                type: 'success',
                removeAfter: 2000,
            })
            return response;
        } catch (error) {
            console.error(error);
            addNotification({
                text: `Could not update node ${node.name}`,
                position: 'top-right',
                type: 'danger',
                removeAfter: 2000,
            })
        }
    };
    async function handleClick(){
        console.log("UPDATING NETWORK", node.is_lighthouse, static_address, groups)
        await updateNode();
        console.log("SEND REQUEST TO UPDATE NODE")
    }
    export const getNodeInfo = async() => {
        try {
            const response = await API.get(`/networks/${netname}/nodes/${nodename}/info`);
            console.log("NODEINFO", response)
            return response;
        } catch (error) {
            console.error(error);
            addNotification({
                text: `Could not get info for node ${node.name}`,
                position: 'top-right',
                type: 'danger',
                removeAfter: 2000,
            })
        }
    };
    async function refreshInfo(){
        const res = await getNodeInfo();
        node = res;
    }
    import { onMount } from 'svelte';
    onMount(async() => {
        await refreshInfo()
    })
    function addGroup(){
        if (newgroup == ""){return}
        console.log("Adding", newgroup, "to", node.groups)
        groups = node.groups
        groups.push(newgroup);
        node.groups = groups;
        newgroup = "";
    }
    function removeGroup(groupname){
        console.log("Removing", groupname, "from", node.groups);
        groups = node.groups;
        for( var i = 0; i < groups.length; i++){ 
            if ( groups[i] === groupname) { 
                groups.splice(i, 1); 
            }
        }
        node.groups = groups;
    }
    function addInboundRule(){
        var rules = node.firewall_inbound
        rules.push({
            port: 'any',
            proto: 'any',
            groups: [],
            any: true,
        })
        node.firewall_inbound = rules
    }
    function addOutboundRule(){
        var rules = node.firewall_outbound
        rules.push({
            port: 'any',
            proto: 'any',
            groups: [],
            any: true,
        })
        node.firewall_outbound = rules
    }
    function deleteInboundRule(index){
        console.log("Deleting rule at ", index)
        var rules = node.firewall_inbound
        rules.splice(index, 1)
        node.firewall_inbound = rules
    }
    function deleteOutboundRule(index){
        console.log("Deleting rule at ", index)
        var rules = node.firewall_outbound
        rules.splice(index, 1)
        node.firewall_outbound = rules
    }
    $: {
        netgroupOptions = netgroups.filter(n => !node.groups.includes(n))
    }
</script>

<main>
    <div class="flex flex-row justify-between my-4">
        <div class="flex flex-row justify-center">
            <label for="lighthouse" class="my-auto">Is Lighthouse</label>
            <input type="checkbox" id="lighthouse" name="lighthouse"
                class="my-auto ml-4"
                bind:checked={node.is_lighthouse}/>
        </div>
        <div class="flex flex-row justify-center">
            <label for="static_address" class="my-auto">Static Address</label>
            <input type="text" id="static_address" name="static_address"
                   class="w-1/2 h-6 pl-3 pr-6 text-base placeholder-gray-600 border rounded-lg appearance-none focus:shadow-outline my-auto mx-4 bg-gray-100"
                   bind:value={node.static_address}/>
        </div>
        <div class="flex flex-row justify-center">
            <label for="listen_port" class="my-auto">Listen Port</label>
            <input type="number" id="listen_port" name="listen_port"
                   class="w-1/4 h-6 pl-3 pr-6 text-base placeholder-gray-600 border rounded-lg appearance-none focus:shadow-outline my-auto mx-4 bg-gray-100"
                   bind:value={node.listen_port} />
        </div>
    </div>
    <div id="groups" class="bg-gray-200 rounded p-4">
        <h1 class="text-xl">Groups</h1>
        <select name="addgroup"
                class="w-1/4 h-10 mt-2 pl-3 pr-6 text-base placeholder-gray-600 border rounded-lg appearance-none focus:shadow-outline my-auto"
                bind:value={newgroup}>
            <option value="">Select group</option>
            {#each netgroupOptions as group}
                <option value={group}>{group}</option>
            {/each}
        </select>
        <button on:click={() => addGroup(newgroup)}
                class="bg-indigo-500 text-white rounded-xl p-1 px-2 h-10"
        >Add Group</button>
        <div class="flex flex-row flex-wrap pt-4 ml-2">
            {#each node.groups as group}
                <div class="bg-purple-400 rounded-xl text-white m-2 p-1 my-auto">
                    {group}
                    <button on:click={() => removeGroup(group)} class="outline-none border-none text-gray-500 mx-2 my-auto">x</button>
                </div>
            {/each}
        </div>
    </div>
    <div id="firewall" class="mt-4">
        <h1 class="text-xl">Firewall Rules</h1>
        <div>
            <h2 class="text-lg">Outbound</h2>
            {#each node.firewall_outbound as fr, i}
                <FirewallRule rule={fr} netgroups={netgroups} index={i}
                    on:deleteRule={() => deleteOutboundRule(i)}/>
            {/each}
            <div class="flex justify-end mt-3">
                <button class="bg-blue-400 rounded-xl text-white p-1"
                    on:click={addOutboundRule}>New</button>
            </div>
        </div> 
        <div>
            <h2 class="text-lg">Inbound</h2>
            {#each node.firewall_inbound as fr, i}
                <FirewallRule rule={fr} netgroups={netgroups} index={i}
                    on:deleteRule={() => deleteInboundRule(i)}/>
            {/each}
            <div class="flex justify-end mt-3">
                <button class="bg-blue-400 rounded-xl text-white p-1"
                    on:click={addInboundRule}>New</button>
            </div>
        </div>
    </div>
    <div id="update">
        <button class="bg-green-400 rounded-xl text-white p-1"
            on:click={async() => await handleClick()}>Update</button>
    </div>
</main>
