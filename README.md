Weekly Review:

	- Process Inbox:
		- Create standalone tasks
		- Create projects
			- Create Project:
				- Define outcome
				- Define next action(s)
	- Identify next actions across all projects and sort by context


Sweep projects for next actions (log projects that don't have any next action):

	gtd tasks > tasks.md

List all projects (display <h1> header, or filename if not present):

	gtd project list

Specify each relevant section from CLI:

	gtd project create -name "Hi" -outcome "Something" -next-action "Something simple" -info "Even more stuff"

Open editor with template as starting point (truncate existing matching project?):

	gtd project create -name "Hi"

Create blank, skip editing (not recommended):

	gtd project create -name "Hi" -blank