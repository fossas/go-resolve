// Code generated by SQLBoiler (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package dbmodels

import (
	"bytes"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"github.com/volatiletech/sqlboiler/strmangle"
	"gopkg.in/volatiletech/null.v6"
)

// Job is an object representing the database table.
type Job struct {
	JobID      int       `boil:"job_id" json:"job_id" toml:"job_id" yaml:"job_id"`
	FaktoryJid string    `boil:"faktory_jid" json:"faktory_jid" toml:"faktory_jid" yaml:"faktory_jid"`
	Started    time.Time `boil:"started" json:"started" toml:"started" yaml:"started"`
	Completed  null.Time `boil:"completed" json:"completed,omitempty" toml:"completed" yaml:"completed,omitempty"`

	R *jobR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L jobL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var JobColumns = struct {
	JobID      string
	FaktoryJid string
	Started    string
	Completed  string
}{
	JobID:      "job_id",
	FaktoryJid: "faktory_jid",
	Started:    "started",
	Completed:  "completed",
}

// jobR is where relationships are stored.
type jobR struct {
	Logs     LogSlice
	Packages PackageSlice
}

// jobL is where Load methods for each relationship are stored.
type jobL struct{}

var (
	jobColumns               = []string{"job_id", "faktory_jid", "started", "completed"}
	jobColumnsWithoutDefault = []string{"faktory_jid", "started", "completed"}
	jobColumnsWithDefault    = []string{"job_id"}
	jobPrimaryKeyColumns     = []string{"job_id"}
)

type (
	// JobSlice is an alias for a slice of pointers to Job.
	// This should generally be used opposed to []Job.
	JobSlice []*Job
	// JobHook is the signature for custom Job hook methods
	JobHook func(boil.Executor, *Job) error

	jobQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	jobType                 = reflect.TypeOf(&Job{})
	jobMapping              = queries.MakeStructMapping(jobType)
	jobPrimaryKeyMapping, _ = queries.BindMapping(jobType, jobMapping, jobPrimaryKeyColumns)
	jobInsertCacheMut       sync.RWMutex
	jobInsertCache          = make(map[string]insertCache)
	jobUpdateCacheMut       sync.RWMutex
	jobUpdateCache          = make(map[string]updateCache)
	jobUpsertCacheMut       sync.RWMutex
	jobUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force bytes in case of primary key column that uses []byte (for relationship compares)
	_ = bytes.MinRead
)
var jobBeforeInsertHooks []JobHook
var jobBeforeUpdateHooks []JobHook
var jobBeforeDeleteHooks []JobHook
var jobBeforeUpsertHooks []JobHook

var jobAfterInsertHooks []JobHook
var jobAfterSelectHooks []JobHook
var jobAfterUpdateHooks []JobHook
var jobAfterDeleteHooks []JobHook
var jobAfterUpsertHooks []JobHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Job) doBeforeInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range jobBeforeInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Job) doBeforeUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range jobBeforeUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Job) doBeforeDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range jobBeforeDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Job) doBeforeUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range jobBeforeUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Job) doAfterInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range jobAfterInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Job) doAfterSelectHooks(exec boil.Executor) (err error) {
	for _, hook := range jobAfterSelectHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Job) doAfterUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range jobAfterUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Job) doAfterDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range jobAfterDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Job) doAfterUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range jobAfterUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddJobHook registers your hook function for all future operations.
func AddJobHook(hookPoint boil.HookPoint, jobHook JobHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		jobBeforeInsertHooks = append(jobBeforeInsertHooks, jobHook)
	case boil.BeforeUpdateHook:
		jobBeforeUpdateHooks = append(jobBeforeUpdateHooks, jobHook)
	case boil.BeforeDeleteHook:
		jobBeforeDeleteHooks = append(jobBeforeDeleteHooks, jobHook)
	case boil.BeforeUpsertHook:
		jobBeforeUpsertHooks = append(jobBeforeUpsertHooks, jobHook)
	case boil.AfterInsertHook:
		jobAfterInsertHooks = append(jobAfterInsertHooks, jobHook)
	case boil.AfterSelectHook:
		jobAfterSelectHooks = append(jobAfterSelectHooks, jobHook)
	case boil.AfterUpdateHook:
		jobAfterUpdateHooks = append(jobAfterUpdateHooks, jobHook)
	case boil.AfterDeleteHook:
		jobAfterDeleteHooks = append(jobAfterDeleteHooks, jobHook)
	case boil.AfterUpsertHook:
		jobAfterUpsertHooks = append(jobAfterUpsertHooks, jobHook)
	}
}

// OneP returns a single job record from the query, and panics on error.
func (q jobQuery) OneP() *Job {
	o, err := q.One()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// One returns a single job record from the query.
func (q jobQuery) One() (*Job, error) {
	o := &Job{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "dbmodels: failed to execute a one query for jobs")
	}

	if err := o.doAfterSelectHooks(queries.GetExecutor(q.Query)); err != nil {
		return o, err
	}

	return o, nil
}

// AllP returns all Job records from the query, and panics on error.
func (q jobQuery) AllP() JobSlice {
	o, err := q.All()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// All returns all Job records from the query.
func (q jobQuery) All() (JobSlice, error) {
	var o []*Job

	err := q.Bind(&o)
	if err != nil {
		return nil, errors.Wrap(err, "dbmodels: failed to assign all query results to Job slice")
	}

	if len(jobAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(queries.GetExecutor(q.Query)); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// CountP returns the count of all Job records in the query, and panics on error.
func (q jobQuery) CountP() int64 {
	c, err := q.Count()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
}

// Count returns the count of all Job records in the query.
func (q jobQuery) Count() (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "dbmodels: failed to count jobs rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table, and panics on error.
func (q jobQuery) ExistsP() bool {
	e, err := q.Exists()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// Exists checks if the row exists in the table.
func (q jobQuery) Exists() (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "dbmodels: failed to check if jobs exists")
	}

	return count > 0, nil
}

// LogsG retrieves all the log's logs.
func (o *Job) LogsG(mods ...qm.QueryMod) logQuery {
	return o.Logs(boil.GetDB(), mods...)
}

// Logs retrieves all the log's logs with an executor.
func (o *Job) Logs(exec boil.Executor, mods ...qm.QueryMod) logQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"logs\".\"job_id\"=?", o.JobID),
	)

	query := Logs(exec, queryMods...)
	queries.SetFrom(query.Query, "\"logs\"")

	if len(queries.GetSelect(query.Query)) == 0 {
		queries.SetSelect(query.Query, []string{"\"logs\".*"})
	}

	return query
}

// PackagesG retrieves all the package's packages.
func (o *Job) PackagesG(mods ...qm.QueryMod) packageQuery {
	return o.Packages(boil.GetDB(), mods...)
}

// Packages retrieves all the package's packages with an executor.
func (o *Job) Packages(exec boil.Executor, mods ...qm.QueryMod) packageQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"packages\".\"job_id\"=?", o.JobID),
	)

	query := Packages(exec, queryMods...)
	queries.SetFrom(query.Query, "\"packages\"")

	if len(queries.GetSelect(query.Query)) == 0 {
		queries.SetSelect(query.Query, []string{"\"packages\".*"})
	}

	return query
}

// LoadLogs allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (jobL) LoadLogs(e boil.Executor, singular bool, maybeJob interface{}) error {
	var slice []*Job
	var object *Job

	count := 1
	if singular {
		object = maybeJob.(*Job)
	} else {
		slice = *maybeJob.(*[]*Job)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &jobR{}
		}
		args[0] = object.JobID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &jobR{}
			}
			args[i] = obj.JobID
		}
	}

	query := fmt.Sprintf(
		"select * from \"logs\" where \"job_id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)
	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load logs")
	}
	defer results.Close()

	var resultSlice []*Log
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice logs")
	}

	if len(logAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(e); err != nil {
				return err
			}
		}
	}
	if singular {
		object.R.Logs = resultSlice
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.JobID == foreign.JobID {
				local.R.Logs = append(local.R.Logs, foreign)
				break
			}
		}
	}

	return nil
}

// LoadPackages allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (jobL) LoadPackages(e boil.Executor, singular bool, maybeJob interface{}) error {
	var slice []*Job
	var object *Job

	count := 1
	if singular {
		object = maybeJob.(*Job)
	} else {
		slice = *maybeJob.(*[]*Job)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &jobR{}
		}
		args[0] = object.JobID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &jobR{}
			}
			args[i] = obj.JobID
		}
	}

	query := fmt.Sprintf(
		"select * from \"packages\" where \"job_id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)
	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load packages")
	}
	defer results.Close()

	var resultSlice []*Package
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice packages")
	}

	if len(packageAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(e); err != nil {
				return err
			}
		}
	}
	if singular {
		object.R.Packages = resultSlice
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.JobID == foreign.JobID {
				local.R.Packages = append(local.R.Packages, foreign)
				break
			}
		}
	}

	return nil
}

// AddLogsG adds the given related objects to the existing relationships
// of the job, optionally inserting them as new records.
// Appends related to o.R.Logs.
// Sets related.R.Job appropriately.
// Uses the global database handle.
func (o *Job) AddLogsG(insert bool, related ...*Log) error {
	return o.AddLogs(boil.GetDB(), insert, related...)
}

// AddLogsP adds the given related objects to the existing relationships
// of the job, optionally inserting them as new records.
// Appends related to o.R.Logs.
// Sets related.R.Job appropriately.
// Panics on error.
func (o *Job) AddLogsP(exec boil.Executor, insert bool, related ...*Log) {
	if err := o.AddLogs(exec, insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddLogsGP adds the given related objects to the existing relationships
// of the job, optionally inserting them as new records.
// Appends related to o.R.Logs.
// Sets related.R.Job appropriately.
// Uses the global database handle and panics on error.
func (o *Job) AddLogsGP(insert bool, related ...*Log) {
	if err := o.AddLogs(boil.GetDB(), insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddLogs adds the given related objects to the existing relationships
// of the job, optionally inserting them as new records.
// Appends related to o.R.Logs.
// Sets related.R.Job appropriately.
func (o *Job) AddLogs(exec boil.Executor, insert bool, related ...*Log) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.JobID = o.JobID
			if err = rel.Insert(exec); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"logs\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"job_id"}),
				strmangle.WhereClause("\"", "\"", 2, logPrimaryKeyColumns),
			)
			values := []interface{}{o.JobID, rel.LogID}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}

			if _, err = exec.Exec(updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.JobID = o.JobID
		}
	}

	if o.R == nil {
		o.R = &jobR{
			Logs: related,
		}
	} else {
		o.R.Logs = append(o.R.Logs, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &logR{
				Job: o,
			}
		} else {
			rel.R.Job = o
		}
	}
	return nil
}

// AddPackagesG adds the given related objects to the existing relationships
// of the job, optionally inserting them as new records.
// Appends related to o.R.Packages.
// Sets related.R.Job appropriately.
// Uses the global database handle.
func (o *Job) AddPackagesG(insert bool, related ...*Package) error {
	return o.AddPackages(boil.GetDB(), insert, related...)
}

// AddPackagesP adds the given related objects to the existing relationships
// of the job, optionally inserting them as new records.
// Appends related to o.R.Packages.
// Sets related.R.Job appropriately.
// Panics on error.
func (o *Job) AddPackagesP(exec boil.Executor, insert bool, related ...*Package) {
	if err := o.AddPackages(exec, insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddPackagesGP adds the given related objects to the existing relationships
// of the job, optionally inserting them as new records.
// Appends related to o.R.Packages.
// Sets related.R.Job appropriately.
// Uses the global database handle and panics on error.
func (o *Job) AddPackagesGP(insert bool, related ...*Package) {
	if err := o.AddPackages(boil.GetDB(), insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddPackages adds the given related objects to the existing relationships
// of the job, optionally inserting them as new records.
// Appends related to o.R.Packages.
// Sets related.R.Job appropriately.
func (o *Job) AddPackages(exec boil.Executor, insert bool, related ...*Package) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.JobID = o.JobID
			if err = rel.Insert(exec); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"packages\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"job_id"}),
				strmangle.WhereClause("\"", "\"", 2, packagePrimaryKeyColumns),
			)
			values := []interface{}{o.JobID, rel.PackageID}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}

			if _, err = exec.Exec(updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.JobID = o.JobID
		}
	}

	if o.R == nil {
		o.R = &jobR{
			Packages: related,
		}
	} else {
		o.R.Packages = append(o.R.Packages, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &packageR{
				Job: o,
			}
		} else {
			rel.R.Job = o
		}
	}
	return nil
}

// JobsG retrieves all records.
func JobsG(mods ...qm.QueryMod) jobQuery {
	return Jobs(boil.GetDB(), mods...)
}

// Jobs retrieves all the records using an executor.
func Jobs(exec boil.Executor, mods ...qm.QueryMod) jobQuery {
	mods = append(mods, qm.From("\"jobs\""))
	return jobQuery{NewQuery(exec, mods...)}
}

// FindJobG retrieves a single record by ID.
func FindJobG(jobID int, selectCols ...string) (*Job, error) {
	return FindJob(boil.GetDB(), jobID, selectCols...)
}

// FindJobGP retrieves a single record by ID, and panics on error.
func FindJobGP(jobID int, selectCols ...string) *Job {
	retobj, err := FindJob(boil.GetDB(), jobID, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// FindJob retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindJob(exec boil.Executor, jobID int, selectCols ...string) (*Job, error) {
	jobObj := &Job{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"jobs\" where \"job_id\"=$1", sel,
	)

	q := queries.Raw(exec, query, jobID)

	err := q.Bind(jobObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "dbmodels: unable to select from jobs")
	}

	return jobObj, nil
}

// FindJobP retrieves a single record by ID with an executor, and panics on error.
func FindJobP(exec boil.Executor, jobID int, selectCols ...string) *Job {
	retobj, err := FindJob(exec, jobID, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *Job) InsertG(whitelist ...string) error {
	return o.Insert(boil.GetDB(), whitelist...)
}

// InsertGP a single record, and panics on error. See Insert for whitelist
// behavior description.
func (o *Job) InsertGP(whitelist ...string) {
	if err := o.Insert(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// InsertP a single record using an executor, and panics on error. See Insert
// for whitelist behavior description.
func (o *Job) InsertP(exec boil.Executor, whitelist ...string) {
	if err := o.Insert(exec, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Insert a single record using an executor.
// Whitelist behavior: If a whitelist is provided, only those columns supplied are inserted
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns without a default value are included (i.e. name, age)
// - All columns with a default, but non-zero are included (i.e. health = 75)
func (o *Job) Insert(exec boil.Executor, whitelist ...string) error {
	if o == nil {
		return errors.New("dbmodels: no jobs provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(jobColumnsWithDefault, o)

	key := makeCacheKey(whitelist, nzDefaults)
	jobInsertCacheMut.RLock()
	cache, cached := jobInsertCache[key]
	jobInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := strmangle.InsertColumnSet(
			jobColumns,
			jobColumnsWithDefault,
			jobColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		cache.valueMapping, err = queries.BindMapping(jobType, jobMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(jobType, jobMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"jobs\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.IndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"jobs\" DEFAULT VALUES"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		if len(wl) != 0 {
			cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRow(cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.Exec(cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "dbmodels: unable to insert into jobs")
	}

	if !cached {
		jobInsertCacheMut.Lock()
		jobInsertCache[key] = cache
		jobInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(exec)
}

// UpdateG a single Job record. See Update for
// whitelist behavior description.
func (o *Job) UpdateG(whitelist ...string) error {
	return o.Update(boil.GetDB(), whitelist...)
}

// UpdateGP a single Job record.
// UpdateGP takes a whitelist of column names that should be updated.
// Panics on error. See Update for whitelist behavior description.
func (o *Job) UpdateGP(whitelist ...string) {
	if err := o.Update(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateP uses an executor to update the Job, and panics on error.
// See Update for whitelist behavior description.
func (o *Job) UpdateP(exec boil.Executor, whitelist ...string) {
	err := o.Update(exec, whitelist...)
	if err != nil {
		panic(boil.WrapErr(err))
	}
}

// Update uses an executor to update the Job.
// Whitelist behavior: If a whitelist is provided, only the columns given are updated.
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns are inferred to start with
// - All primary keys are subtracted from this set
// Update does not automatically update the record in case of default values. Use .Reload()
// to refresh the records.
func (o *Job) Update(exec boil.Executor, whitelist ...string) error {
	var err error
	if err = o.doBeforeUpdateHooks(exec); err != nil {
		return err
	}
	key := makeCacheKey(whitelist, nil)
	jobUpdateCacheMut.RLock()
	cache, cached := jobUpdateCache[key]
	jobUpdateCacheMut.RUnlock()

	if !cached {
		wl := strmangle.UpdateColumnSet(
			jobColumns,
			jobPrimaryKeyColumns,
			whitelist,
		)

		if len(whitelist) == 0 {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return errors.New("dbmodels: unable to update jobs, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"jobs\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, jobPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(jobType, jobMapping, append(wl, jobPrimaryKeyColumns...))
		if err != nil {
			return err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	_, err = exec.Exec(cache.query, values...)
	if err != nil {
		return errors.Wrap(err, "dbmodels: unable to update jobs row")
	}

	if !cached {
		jobUpdateCacheMut.Lock()
		jobUpdateCache[key] = cache
		jobUpdateCacheMut.Unlock()
	}

	return o.doAfterUpdateHooks(exec)
}

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q jobQuery) UpdateAllP(cols M) {
	if err := q.UpdateAll(cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values.
func (q jobQuery) UpdateAll(cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "dbmodels: unable to update all for jobs")
	}

	return nil
}

// UpdateAllG updates all rows with the specified column values.
func (o JobSlice) UpdateAllG(cols M) error {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAllGP updates all rows with the specified column values, and panics on error.
func (o JobSlice) UpdateAllGP(cols M) {
	if err := o.UpdateAll(boil.GetDB(), cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAllP updates all rows with the specified column values, and panics on error.
func (o JobSlice) UpdateAllP(exec boil.Executor, cols M) {
	if err := o.UpdateAll(exec, cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o JobSlice) UpdateAll(exec boil.Executor, cols M) error {
	ln := int64(len(o))
	if ln == 0 {
		return nil
	}

	if len(cols) == 0 {
		return errors.New("dbmodels: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), jobPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"jobs\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, jobPrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "dbmodels: unable to update all in job slice")
	}

	return nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *Job) UpsertG(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	return o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...)
}

// UpsertGP attempts an insert, and does an update or ignore on conflict. Panics on error.
func (o *Job) UpsertGP(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *Job) UpsertP(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(exec, updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *Job) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	if o == nil {
		return errors.New("dbmodels: no jobs provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(jobColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs postgres problems
	buf := strmangle.GetBuffer()

	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range updateColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range whitelist {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	jobUpsertCacheMut.RLock()
	cache, cached := jobUpsertCache[key]
	jobUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := strmangle.InsertColumnSet(
			jobColumns,
			jobColumnsWithDefault,
			jobColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		update := strmangle.UpdateColumnSet(
			jobColumns,
			jobPrimaryKeyColumns,
			updateColumns,
		)
		if len(update) == 0 {
			return errors.New("dbmodels: unable to upsert jobs, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(jobPrimaryKeyColumns))
			copy(conflict, jobPrimaryKeyColumns)
		}
		cache.query = queries.BuildUpsertQueryPostgres(dialect, "\"jobs\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(jobType, jobMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(jobType, jobMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRow(cache.query, vals...).Scan(returns...)
		if err == sql.ErrNoRows {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.Exec(cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "dbmodels: unable to upsert jobs")
	}

	if !cached {
		jobUpsertCacheMut.Lock()
		jobUpsertCache[key] = cache
		jobUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(exec)
}

// DeleteP deletes a single Job record with an executor.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *Job) DeleteP(exec boil.Executor) {
	if err := o.Delete(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteG deletes a single Job record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *Job) DeleteG() error {
	if o == nil {
		return errors.New("dbmodels: no Job provided for deletion")
	}

	return o.Delete(boil.GetDB())
}

// DeleteGP deletes a single Job record.
// DeleteGP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *Job) DeleteGP() {
	if err := o.DeleteG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Delete deletes a single Job record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Job) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("dbmodels: no Job provided for delete")
	}

	if err := o.doBeforeDeleteHooks(exec); err != nil {
		return err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), jobPrimaryKeyMapping)
	sql := "DELETE FROM \"jobs\" WHERE \"job_id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "dbmodels: unable to delete from jobs")
	}

	if err := o.doAfterDeleteHooks(exec); err != nil {
		return err
	}

	return nil
}

// DeleteAllP deletes all rows, and panics on error.
func (q jobQuery) DeleteAllP() {
	if err := q.DeleteAll(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all matching rows.
func (q jobQuery) DeleteAll() error {
	if q.Query == nil {
		return errors.New("dbmodels: no jobQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "dbmodels: unable to delete all from jobs")
	}

	return nil
}

// DeleteAllGP deletes all rows in the slice, and panics on error.
func (o JobSlice) DeleteAllGP() {
	if err := o.DeleteAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAllG deletes all rows in the slice.
func (o JobSlice) DeleteAllG() error {
	if o == nil {
		return errors.New("dbmodels: no Job slice provided for delete all")
	}
	return o.DeleteAll(boil.GetDB())
}

// DeleteAllP deletes all rows in the slice, using an executor, and panics on error.
func (o JobSlice) DeleteAllP(exec boil.Executor) {
	if err := o.DeleteAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o JobSlice) DeleteAll(exec boil.Executor) error {
	if o == nil {
		return errors.New("dbmodels: no Job slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	if len(jobBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), jobPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"jobs\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, jobPrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "dbmodels: unable to delete all from job slice")
	}

	if len(jobAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	return nil
}

// ReloadGP refetches the object from the database and panics on error.
func (o *Job) ReloadGP() {
	if err := o.ReloadG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadP refetches the object from the database with an executor. Panics on error.
func (o *Job) ReloadP(exec boil.Executor) {
	if err := o.Reload(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadG refetches the object from the database using the primary keys.
func (o *Job) ReloadG() error {
	if o == nil {
		return errors.New("dbmodels: no Job provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Job) Reload(exec boil.Executor) error {
	ret, err := FindJob(exec, o.JobID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllGP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *JobSlice) ReloadAllGP() {
	if err := o.ReloadAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *JobSlice) ReloadAllP(exec boil.Executor) {
	if err := o.ReloadAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *JobSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("dbmodels: empty JobSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *JobSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	jobs := JobSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), jobPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"jobs\".* FROM \"jobs\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, jobPrimaryKeyColumns, len(*o))

	q := queries.Raw(exec, sql, args...)

	err := q.Bind(&jobs)
	if err != nil {
		return errors.Wrap(err, "dbmodels: unable to reload all in JobSlice")
	}

	*o = jobs

	return nil
}

// JobExists checks if the Job row exists.
func JobExists(exec boil.Executor, jobID int) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"jobs\" where \"job_id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, jobID)
	}

	row := exec.QueryRow(sql, jobID)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "dbmodels: unable to check if jobs exists")
	}

	return exists, nil
}

// JobExistsG checks if the Job row exists.
func JobExistsG(jobID int) (bool, error) {
	return JobExists(boil.GetDB(), jobID)
}

// JobExistsGP checks if the Job row exists. Panics on error.
func JobExistsGP(jobID int) bool {
	e, err := JobExists(boil.GetDB(), jobID)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// JobExistsP checks if the Job row exists. Panics on error.
func JobExistsP(exec boil.Executor, jobID int) bool {
	e, err := JobExists(exec, jobID)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}
