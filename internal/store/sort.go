package store

import (
	"context"
)

// SortOrder specifies how to sort catalog items.
type SortOrder string

const (
	SortNameAsc    SortOrder = "name_asc"    // A-Z
	SortNameDesc   SortOrder = "name_desc"   // Z-A
	SortRatingDesc SortOrder = "rating_desc" // Highest rating first
	SortRatingAsc  SortOrder = "rating_asc"  // Lowest rating first
	SortRecent     SortOrder = "recent"      // Newest added first
	SortOld        SortOrder = "old"         // Oldest added first
)

// ListVODsSorted returns VODs in a category sorted by the specified order.
func (s *Store) ListVODsSorted(ctx context.Context, categoryID int, order SortOrder) ([]VODRow, error) {
	var query string
	switch order {
	case SortNameDesc:
		query = `SELECT stream_id,category_id,name,COALESCE(year,0),COALESCE(plot,''),
		                COALESCE(stream_icon_url,''),COALESCE(container_extension,''),
		                COALESCE(added,0),COALESCE(rating,0)
		         FROM vods WHERE category_id=? ORDER BY name DESC`
	case SortRatingDesc:
		query = `SELECT stream_id,category_id,name,COALESCE(year,0),COALESCE(plot,''),
		                COALESCE(stream_icon_url,''),COALESCE(container_extension,''),
		                COALESCE(added,0),COALESCE(rating,0)
		         FROM vods WHERE category_id=? ORDER BY COALESCE(rating,0) DESC, name ASC`
	case SortRatingAsc:
		query = `SELECT stream_id,category_id,name,COALESCE(year,0),COALESCE(plot,''),
		                COALESCE(stream_icon_url,''),COALESCE(container_extension,''),
		                COALESCE(added,0),COALESCE(rating,0)
		         FROM vods WHERE category_id=? ORDER BY COALESCE(rating,0) ASC, name ASC`
	case SortRecent:
		query = `SELECT stream_id,category_id,name,COALESCE(year,0),COALESCE(plot,''),
		                COALESCE(stream_icon_url,''),COALESCE(container_extension,''),
		                COALESCE(added,0),COALESCE(rating,0)
		         FROM vods WHERE category_id=? ORDER BY COALESCE(added,0) DESC, name ASC`
	case SortOld:
		query = `SELECT stream_id,category_id,name,COALESCE(year,0),COALESCE(plot,''),
		                COALESCE(stream_icon_url,''),COALESCE(container_extension,''),
		                COALESCE(added,0),COALESCE(rating,0)
		         FROM vods WHERE category_id=? ORDER BY COALESCE(added,0) ASC, name ASC`
	default: // SortNameAsc
		query = `SELECT stream_id,category_id,name,COALESCE(year,0),COALESCE(plot,''),
		                COALESCE(stream_icon_url,''),COALESCE(container_extension,''),
		                COALESCE(added,0),COALESCE(rating,0)
		         FROM vods WHERE category_id=? ORDER BY name ASC`
	}

	rows, err := s.db.QueryContext(ctx, query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []VODRow
	for rows.Next() {
		var v VODRow
		if err := rows.Scan(&v.StreamID, &v.CategoryID, &v.Name, &v.Year, &v.Plot,
			&v.StreamIcon, &v.ContainerExt, &v.Added, &v.Rating); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}

// ListSeriesSorted returns series in a category sorted by the specified order.
func (s *Store) ListSeriesSorted(ctx context.Context, categoryID int, order SortOrder) ([]SeriesRow, error) {
	var query string
	switch order {
	case SortNameDesc:
		query = `SELECT series_id,category_id,name,COALESCE(year,0),COALESCE(plot,''),
		                COALESCE(cover_url,''),COALESCE(backdrop_url,'')
		         FROM series WHERE category_id=? ORDER BY name DESC`
	case SortRatingDesc, SortRatingAsc:
		// Series don't have ratings, fall back to name sorting
		query = `SELECT series_id,category_id,name,COALESCE(year,0),COALESCE(plot,''),
		                COALESCE(cover_url,''),COALESCE(backdrop_url,'')
		         FROM series WHERE category_id=? ORDER BY name ASC`
	case SortRecent:
		query = `SELECT series_id,category_id,name,COALESCE(year,0),COALESCE(plot,''),
		                COALESCE(cover_url,''),COALESCE(backdrop_url,'')
		         FROM series WHERE category_id=? ORDER BY series_id DESC, name ASC`
	case SortOld:
		query = `SELECT series_id,category_id,name,COALESCE(year,0),COALESCE(plot,''),
		                COALESCE(cover_url,''),COALESCE(backdrop_url,'')
		         FROM series WHERE category_id=? ORDER BY series_id ASC, name ASC`
	default: // SortNameAsc
		query = `SELECT series_id,category_id,name,COALESCE(year,0),COALESCE(plot,''),
		                COALESCE(cover_url,''),COALESCE(backdrop_url,'')
		         FROM series WHERE category_id=? ORDER BY name ASC`
	}

	rows, err := s.db.QueryContext(ctx, query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []SeriesRow
	for rows.Next() {
		var r SeriesRow
		if err := rows.Scan(&r.SeriesID, &r.CategoryID, &r.Name, &r.Year, &r.Plot,
			&r.CoverURL, &r.BackdropURL); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

// ListAllVODsSorted returns all VODs sorted by the specified order.
func (s *Store) ListAllVODsSorted(ctx context.Context, order SortOrder) ([]VODRow, error) {
	var query string
	switch order {
	case SortNameDesc:
		query = `SELECT stream_id,category_id,name,COALESCE(year,0),COALESCE(plot,''),
		                COALESCE(stream_icon_url,''),COALESCE(container_extension,''),
		                COALESCE(added,0),COALESCE(rating,0)
		         FROM vods ORDER BY name DESC`
	case SortRatingDesc:
		query = `SELECT stream_id,category_id,name,COALESCE(year,0),COALESCE(plot,''),
		                COALESCE(stream_icon_url,''),COALESCE(container_extension,''),
		                COALESCE(added,0),COALESCE(rating,0)
		         FROM vods ORDER BY COALESCE(rating,0) DESC, name ASC`
	case SortRatingAsc:
		query = `SELECT stream_id,category_id,name,COALESCE(year,0),COALESCE(plot,''),
		                COALESCE(stream_icon_url,''),COALESCE(container_extension,''),
		                COALESCE(added,0),COALESCE(rating,0)
		         FROM vods ORDER BY COALESCE(rating,0) ASC, name ASC`
	case SortRecent:
		query = `SELECT stream_id,category_id,name,COALESCE(year,0),COALESCE(plot,''),
		                COALESCE(stream_icon_url,''),COALESCE(container_extension,''),
		                COALESCE(added,0),COALESCE(rating,0)
		         FROM vods ORDER BY COALESCE(added,0) DESC, name ASC`
	case SortOld:
		query = `SELECT stream_id,category_id,name,COALESCE(year,0),COALESCE(plot,''),
		                COALESCE(stream_icon_url,''),COALESCE(container_extension,''),
		                COALESCE(added,0),COALESCE(rating,0)
		         FROM vods ORDER BY COALESCE(added,0) ASC, name ASC`
	default: // SortNameAsc
		query = `SELECT stream_id,category_id,name,COALESCE(year,0),COALESCE(plot,''),
		                COALESCE(stream_icon_url,''),COALESCE(container_extension,''),
		                COALESCE(added,0),COALESCE(rating,0)
		         FROM vods ORDER BY name ASC`
	}

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []VODRow
	for rows.Next() {
		var v VODRow
		if err := rows.Scan(&v.StreamID, &v.CategoryID, &v.Name, &v.Year, &v.Plot,
			&v.StreamIcon, &v.ContainerExt, &v.Added, &v.Rating); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}

// ListAllSeriesSorted returns all series sorted by the specified order.
func (s *Store) ListAllSeriesSorted(ctx context.Context, order SortOrder) ([]SeriesRow, error) {
	var query string
	switch order {
	case SortNameDesc:
		query = `SELECT series_id,category_id,name,COALESCE(year,0),COALESCE(plot,''),
		                COALESCE(cover_url,''),COALESCE(backdrop_url,'')
		         FROM series ORDER BY name DESC`
	case SortRatingDesc, SortRatingAsc:
		// Series don't have ratings, fall back to name sorting
		query = `SELECT series_id,category_id,name,COALESCE(year,0),COALESCE(plot,''),
		                COALESCE(cover_url,''),COALESCE(backdrop_url,'')
		         FROM series ORDER BY name ASC`
	case SortRecent:
		query = `SELECT series_id,category_id,name,COALESCE(year,0),COALESCE(plot,''),
		                COALESCE(cover_url,''),COALESCE(backdrop_url,'')
		         FROM series ORDER BY series_id DESC, name ASC`
	case SortOld:
		query = `SELECT series_id,category_id,name,COALESCE(year,0),COALESCE(plot,''),
		                COALESCE(cover_url,''),COALESCE(backdrop_url,'')
		         FROM series ORDER BY series_id ASC, name ASC`
	default: // SortNameAsc
		query = `SELECT series_id,category_id,name,COALESCE(year,0),COALESCE(plot,''),
		                COALESCE(cover_url,''),COALESCE(backdrop_url,'')
		         FROM series ORDER BY name ASC`
	}

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []SeriesRow
	for rows.Next() {
		var r SeriesRow
		if err := rows.Scan(&r.SeriesID, &r.CategoryID, &r.Name, &r.Year, &r.Plot,
			&r.CoverURL, &r.BackdropURL); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}
